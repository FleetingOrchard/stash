import React from "react";
import { Helmet } from "react-helmet";
import { useIntl } from "react-intl";
import { Route, Switch } from "react-router-dom";
import { TITLE_SUFFIX } from "../Shared";
import { BookmarkList } from "./BookmarkList";

const Bookmarks: React.FC = () => {
  const intl = useIntl();

  const title_template = `${intl.formatMessage({
    id: "bookmarks",
  })} ${TITLE_SUFFIX}`;
  return (
    <>
      <Helmet
        defaultTitle={title_template}
        titleTemplate={`%s | ${title_template}`}
      />

      <Switch>
        <Route exact path="/bookmarks" component={BookmarkList} />
      </Switch>
    </>
  );
};

export default Bookmarks;