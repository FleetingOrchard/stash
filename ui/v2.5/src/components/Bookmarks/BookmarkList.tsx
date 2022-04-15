import React, { useEffect, useState } from "react";
import { ReactSortable, SortableEvent, Sortable } from "react-sortablejs"
import { Bookmark } from "./Bookmark";
import { useBookmarkUpdate } from "src/core/StashService"
import * as GQL from "src/core/generated-graphql";

interface IBookmarkList {
}

export const BookmarkList: React.FC<IBookmarkList> = () => {
  const bookmarksQuery = GQL.useAllBookmarksQuery();
  const [updateBookmark] = useBookmarkUpdate();
  const [listState, setListState] = useState(bookmarksQuery.data?.allBookmarks);

  // Reset list on changes to bookmarksQuery.
  useEffect(() => {
    setListState(bookmarksQuery.data?.allBookmarks);
  }, [bookmarksQuery.data?.allBookmarks]);

  const storeGet = (sortable: Sortable) => {
    return sortable.toArray();
  }

  const storeSet = (sortable: Sortable) => {
    if (bookmarksQuery.data === undefined) {
      console.error("BookmarkList: storeSet: bookmarsQuery.data is undefined!");
      return;
    }
    const bms = bookmarksQuery.data.allBookmarks;

    // TODO: Change this to all happen on-server.
    const posList = sortable.toArray(); 
    for (var i = 0; i < bms.length; i++)
    {
      const bm = bms.find(x => x.id === posList[i]);
      if (bm === undefined)
        continue;
      bm.position = i;
    }
    bms.forEach(bm => updateBookmark({variables: {input: {id: bm.id, url: bm.url, name: bm.name, position: bm.position}}}));
  }

  const renderList = () => {
    if (bookmarksQuery.loading || listState === undefined)
    {
      return null; 
    }

    return (
      <div className="row px-xl-5 justify-content-center">
        <ReactSortable className="w-100" list={listState} setList={setListState} store={{get: storeGet, set: storeSet}}>
          {listState?.map((item) => (
            <Bookmark bookmark={item} />
          ))}
        </ReactSortable>
        <Bookmark bookmark={{ id: "", url: "" } as GQL.Bookmark} />
      </div>
    );
  }

  return renderList();
};