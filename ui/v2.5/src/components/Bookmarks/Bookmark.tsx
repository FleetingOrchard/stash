import React, { useEffect, useState, CSSProperties } from "react";
import { Button, FormControl, InputGroup } from "react-bootstrap";
import _ from "lodash";
import { Icon } from "src/components/Shared";
import * as GQL from "src/core/generated-graphql";
import { useBookmarkCreate, useBookmarkDestroy, useBookmarkUpdate } from "src/core/StashService"
import { useToast } from "src/hooks";
import { faTrash, faEdit, faLink, faPlus } from "@fortawesome/free-solid-svg-icons";

interface IProps {
  bookmark: GQL.Bookmark,
}

export const Bookmark: React.FC<IProps> = ({
  bookmark,
}) => {
  const Toast = useToast();
  const [isNew] = useState(bookmark.id === "");
  const [isEditing, setEditing] = useState(bookmark.id === "");
  const [updateNeeded, setUpdateNeed] = useState(false);
  const [btnsActive, setBtnsActive] = useState(true);
  const [stateUrl, setStateUrl] = useState(bookmark.url);
  const [stateName, setStateName] = useState(bookmark.name);
  const [deleteBookmark] = useBookmarkDestroy(bookmark);
  const [createBookmark] = useBookmarkCreate();
  const [updateBookmark] = useBookmarkUpdate();

  useEffect(() => {
    if (!updateNeeded)
    {
      setStateUrl(bookmark.url);
      setStateName(bookmark.name);
    }
  });

  function setUrl(value: string) {
    setUpdateNeed(true);
    setStateUrl(value);
  }

  function setName(value: string) {
    setUpdateNeed(true);
    setStateName(value);
  }

  function toggleEdit() {
    if (isEditing && updateNeeded)
    {
      setBtnsActive(false);
      updateBookmark({variables: {input: {id:bookmark.id, url:stateUrl, name:stateName, position:bookmark.position}}})
      .then(() => { Toast.success({content: "Bookmark edited!"}); })
      .catch((e) => { Toast.error(e);})
      .finally(() => {
        setEditing(false);
        setBtnsActive(true);
        setUpdateNeed(false);
      });
    }
    else
    {
      setEditing(!isEditing);
    }
  }

  function createNewBookmark() {
    setBtnsActive(false);
    createBookmark({variables: {input: {url: stateUrl, name:stateName}}})
      .then(() => { Toast.success({content: "Bookmark created!"}); })
      .catch((e) => { Toast.error(e);})
      .finally(() => {
        setBtnsActive(true);
        setUpdateNeed(false);
        setUrl("");
      });
  }

  function deleteTheBookmark() {
    setBtnsActive(false);
    deleteBookmark({variables: {id: bookmark.id}})
      .then(() => { Toast.success({content: "Bookmark deleted!"}); })
      .catch((e) => { Toast.error(e);})
      .finally(() => { setBtnsActive(true); });
  }

  function renderUrlAndNameButton() {
    const buttonProps: CSSProperties = {
      overflow: "clip",
      wordBreak: "break-all",
    };

    if (isEditing) {
      return [
        <FormControl type="text" value={stateUrl} onChange={e => setUrl(e.target.value)} />,
        <FormControl type="text" value={stateName} onChange={e => setName(e.target.value)} />
      ];
    }

    return [
      <Button style={buttonProps} bsPrefix="form-control" className="text-light bg-secondary border-secondary text-left" href={bookmark.url}>{decodeURI(bookmark.url)}</Button>,
      <Button style={buttonProps} bsPrefix="form-control" className="text-light bg-secondary border-secondary text-left" href={bookmark.url}>{bookmark.name}</Button>
    ];
  }

  function renderTrailingButtons() {
    if (!isNew) {
      return (
        <div className="input-group-append">
          <Button variant="info"
            disabled={!btnsActive}
            onClick={toggleEdit} >
            <Icon icon={faEdit} />
          </Button>
          <Button variant="danger"
            disabled={!btnsActive}
            onClick={deleteTheBookmark} >
            <Icon icon={faTrash} />
          </Button>
        </div>
      )
    }

    return (
      <div className="input-group-append">
        <Button variant="info"
          disabled={!btnsActive}
          onClick={createNewBookmark} >
          <Icon icon={faPlus} />
        </Button>
      </div>
    )
  }

  return (
    <InputGroup className="p-2 w-100" data-id={bookmark.id}>
      <div className="input-group-prepend">
        <Button href={bookmark.url} variant="primary">
          <Icon icon={faLink} />
        </Button>
      </div>
      {renderUrlAndNameButton()}
      {renderTrailingButtons()}
    </InputGroup>
  );
};