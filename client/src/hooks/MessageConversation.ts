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
    selectedConversation: IInboxEntry | undefined,
    chatboxRef: React.RefObject<HTMLDivElement | null>
) => {
    const [loaded, setLoaded] = useState<boolean>(false);
    const [isAtTop, setIsAtTop] = useState<boolean>(true);
    const [preventFetch, setPreventFetch] = useState<boolean>(false);

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
            if (!selectedConversation) return [];
            const data = messageStore.getConversationMessages(
                selectedConversation.conversation_id
            );
            // In case messages already exist, then a new query trigger means that the user wants to fetch older messages, so include this as part of query
            let oldestMsgID: string | undefined = undefined;
            if (data) {
                oldestMsgID = data[data.length - 1].id;
            }

            const resp = await fetchConversationMessages(
                selectedConversation!.conversation_id,
                20,
                oldestMsgID
            );
            if (resp.error) {
                throw new APIError(resp.message, resp.status);
            }

            if (resp.data.length == 0) {
                // If query returns no messages, that means either the user has scrolled to the top of the conversation, or theres no messages in the conversation. No point in fetching again.
                setPreventFetch(true);
            }

            setLoaded(true);
            return resp.data;
        },
        retry: false,
        enabled: false,
    });

    const handleScroll = () => {
        // Handle scroll for fetching older messages
        if (chatboxRef.current) {
            if (
                Math.abs(
                    chatboxRef.current.clientHeight -
                        chatboxRef.current.scrollTop -
                        chatboxRef.current.scrollHeight
                ) < 5
            ) {
                if (!isAtTop) {
                    // Prevent re-triggering if already at top

                    if (!messageQuery.isFetching && !preventFetch) {
                        messageQuery.refetch();
                        setIsAtTop(true);
                        
                    }
                }
            } else {
                if (isAtTop) {
                    setIsAtTop(false); // Update state when no longer at top
                }
            }
        }
    };

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
                setLoaded(true);
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
                        messageQuery.data
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
        return {
            data: messageStore.getConversationMessages(
                selectedConversation.conversation_id
            ),
            handleScroll,
            loaded,
        };
    } else {
        return { data: [], handleScroll, loaded };
    }
};
