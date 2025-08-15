import type React from "react";
import { useEffect, useRef, useState } from "react";
import { isEmptyString, log } from "@/utils";
import { useUIStore } from "@/store/ui";
import { useInboxStore } from "@/store/inbox";
import type { IInboxEntry } from "@/ws/types";
import { InboxEntry } from "../InboxEntry";
import { Input } from "@/components/ui/input";
import { Spinner } from "@/components/ui/shadcn-io/spinner";
import { searchUser } from "@/api/search";

const searchForRelevantUser = async (
  username: string,
  inbox?: IInboxEntry[]
): Promise<IInboxEntry[] | null> => {
  const usernameLower = username.toLowerCase();

  // 1. Try inbox search
  if (inbox && inbox.length > 0) {
    const filtered = inbox.filter((convo) =>
      convo.name.toLowerCase().includes(usernameLower)
    );

    filtered.sort((a, b) => {
      // exact match first
      if (a.name.toLowerCase() === usernameLower) return -1;
      if (b.name.toLowerCase() === usernameLower) return 1;

      // starts with first
      if (a.name.toLowerCase().startsWith(usernameLower)) return -1;
      if (b.name.toLowerCase().startsWith(usernameLower)) return 1;

      // fallback: alphabetical
      return a.name.localeCompare(b.name);
    });

    if (filtered.length > 0) {
      return filtered;
    }
  }

  try {
    const response = await searchUser(username);
    if (response) {
      const user = response.data[0];
      return [
        {
          conversation_id: "",
          name: user.display_name,
          updated_at: 0,
          user_id: user.id,
          type: "DM",
        },
      ];
    }
  } catch (err) {
    log("error", "searching user:", err);
  }

  return [];
};

export const ConversationSearchComponent: React.FC = () => {
  const uiStore = useUIStore();
  const inboxState = useInboxStore();

  const [searchResults, setSearchResults] = useState<IInboxEntry[] | null>(
    null
  );
  const isFirstMount = useRef(true);

  const fetchAsync = async () => {
    const searchTerm = uiStore.currentConvoSearchInput!;
    const results = await searchForRelevantUser(
      searchTerm,
      inboxState.conversations
    );
    setSearchResults(results);
  };

  useEffect(() => {
    isFirstMount.current = false;
  }, []);

  useEffect(() => {
    if (!uiStore.currentConvoSearchInput) return;
    fetchAsync();
  }, [uiStore.currentConvoSearchInput]);

  useEffect(() => {
    setTimeout(async () => {
      if (!searchResults && uiStore.isSearchingConvo) {
        await fetchAsync();
      }
    }, 750);
  }, [searchResults]);

  return (
    <div className="grid gap-2 overflow-y-auto">
      {!searchResults && uiStore.isSearchingConvo ? (
        <div className="w-full h-full flex items-center justify-center">
          <Spinner />
        </div>
      ) : (
        <>
          {searchResults && (
            <>
              {searchResults.length > 0 ? (
                searchResults.map((result) => (
                  <InboxEntry
                    conversationID={result.conversation_id}
                    user_id={result.user_id}
                    last_msg_time={result.updated_at}
                    name={result.name}
                    onClick={() => uiStore.setIsSearchingConvo(false)}
                  />
                ))
              ) : (
                <p className="mx-5 text-center text-foreground/30 text-2xl py-10 font-normal">
                  No users found
                </p>
              )}
            </>
          )}
        </>
      )}
    </div>
  );
};

export const ConversationSearchbar: React.FC = () => {
  const uiStore = useUIStore();

  const [username, setUsername] = useState<string>("");

  useEffect(() => {
    // simple debounce
    const searchConvo = setTimeout(async () => {
      if (!isEmptyString(username)) {
        uiStore.setIsSearchingConvo(true);
        uiStore.setCurrentConvoSearchInput(username);
      } else {
        uiStore.setIsSearchingConvo(false);
        uiStore.setCurrentConvoSearchInput(); // Set to undefined
      }
    }, 750);
    return () => clearTimeout(searchConvo);
  }, [username]);

  return (
    <div className="flex h-14 justify-center w-full">
      <div className="mx-5 flex items-center gap-5 w-full relative">
        <Input
          type="text"
          className="rounded-xl"
          placeholder="Search Usernames"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
      </div>
    </div>
  );
};
