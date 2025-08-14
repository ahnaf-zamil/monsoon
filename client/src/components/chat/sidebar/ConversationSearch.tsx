import type React from "react";
import { useEffect, useRef, useState } from "react";
import { isEmptyString } from "@/utils";
import { useUIStore } from "@/store/ui";
import { useInboxStore } from "@/store/inbox";
import type { IInboxEntry } from "@/ws/types";
import { InboxEntry } from "../InboxEntry";
import { Input } from "@/components/ui/input";

export const ConversationSearchComponent: React.FC = () => {
  const uiStore = useUIStore();
  const inboxState = useInboxStore();

  const [searchResults, setSearchResults] = useState<IInboxEntry[] | null>(
    null
  );
  const isFirstMount = useRef(true);

  useEffect(() => {
    isFirstMount.current = false;
  }, []);

  useEffect(() => {
    if (!uiStore.currentConvoSearchInput) return;
    const searchTerm = uiStore.currentConvoSearchInput!;
    let results = inboxState.conversations.filter((convo) =>
      convo.name.toLowerCase().includes(searchTerm.toLowerCase())
    );

    results = results.sort((a, b) => {
      const term = searchTerm.toLowerCase();

      // names with exact matches are first
      if (a.name.toLowerCase() === term && b.name.toLowerCase() !== term) {
        return -1;
      }
      if (a.name.toLowerCase() !== term && b.name.toLowerCase() === term) {
        return 1;
      }

      // names starting with term
      if (
        a.name.toLowerCase().startsWith(term) &&
        !b.name.toLowerCase().startsWith(term)
      ) {
        return -1;
      }
      if (
        !a.name.toLowerCase().startsWith(term) &&
        b.name.toLowerCase().startsWith(term)
      ) {
        return 1;
      }

      // if none of them work, just sort alphabetically
      return a.name.localeCompare(b.name);
    });

    setSearchResults(results);
  }, [uiStore.currentConvoSearchInput]);

  return (
    <div className="grid gap-2 overflow-y-auto">
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
