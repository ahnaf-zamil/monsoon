import { useQuery } from "@tanstack/react-query";
import { HttpStatusCode } from "axios";
import { useEffect, useState } from "react";
import { fetchConversationMessages } from "../api/message";
import { type IMessageData, APIError } from "../api/types";
import { useCurrentUser } from "../context/AuthContext";
import { useMessageStore } from "../store/message";
import { useUserCacheStore } from "../store/user";
import type { IInboxEntry } from "../ws/types";
import { log } from "../utils";

export const useMessagesForConversation = (
    selectedConversation: IInboxEntry | undefined
) => {
    const [loaded, setLoaded] = useState<boolean>(false);

    const userCache = useUserCacheStore();
    const messageStore = useMessageStore();
    const currentUser = useCurrentUser();

    const fetchAndCacheConversationUser = async (
        _conversationID: string,
        userID: string
    ) => {
        // WIP: Fetch DM user data from API and store in cache
        const cachedUser = userCache.getUser(userID);

        if (!cachedUser) {
            console.log("User data not cached");
        } else {
            console.log("User data is cached");
        }
    };

    const messageQuery = useQuery({
        queryKey: ["chat-messages"],
        queryFn: async (): Promise<IMessageData[]> => {
            const resp = await fetchConversationMessages(
                selectedConversation!.conversation_id,
                20
            );
            if (resp.error) {
                throw new APIError(resp.message, resp.status);
            }
            setLoaded(true)
            return resp.data;
        },
        retry: false,
        enabled: false,
    });

    useEffect(() => {
        if (selectedConversation != null) {
            const convoID = selectedConversation.conversation_id;

            if (selectedConversation.type === "DM") {
                // Cache conversation user, WIP
                fetchAndCacheConversationUser(
                    selectedConversation.conversation_id,
                    selectedConversation.user_id!
                ); // user_id will always be present in DM conversations
            }

            // Check for caching and fetch if cache miss
            const msg = messageStore.getConversationMessages(convoID);
            if (msg == undefined) {
                messageQuery.refetch();
            } else {
                setLoaded(true)
            }
        }
    }, [selectedConversation]);

    useEffect(() => {
        // Store fetched messages into cache
        if (selectedConversation?.conversation_id && messageQuery.isSuccess) {
            if (messageQuery.data.length > 0) {
                const firstMsg = messageQuery.data[0];
                if (
                    firstMsg.conversation_id ===
                    selectedConversation.conversation_id
                ) {
                    messageStore.storeMessages(
                        selectedConversation?.conversation_id,
                        messageQuery.data!
                    );
                } else {
                    log(
                        "warn",
                        "Mismatched conversation data — skipping store"
                    );
                }
            }
        }
    }, [messageQuery.data]); // Not using selectedConversation?.conversation_id in deps to prevent incorrect state sync. Just because selected convo changes doesn't mean query has new data, since sometimes the data may just be cached and query wont run

    useEffect(() => {
        const retryAfterAuth = async () => {
            // If access token has expired, then refetch token and fetch messages again
            await currentUser.refetchUser();
            if (currentUser.isSuccess && !currentUser.isError) {
                await messageQuery.refetch();
            }
        };

        if (
            messageQuery.isError &&
            messageQuery.error instanceof APIError &&
            messageQuery.error.status === HttpStatusCode.Unauthorized
        ) {
            retryAfterAuth();
        } else if (
            messageQuery.isError &&
            messageQuery.error instanceof APIError
        ) {
            log("error", "Unhandled error status:", messageQuery.error);
        }
    }, [messageQuery.isError, messageQuery.error]);

    if (selectedConversation) {
        return {data: messageStore.getConversationMessages(
            selectedConversation.conversation_id
        ), loaded};
    } else {
        return {data: [], loaded};
    }
};
