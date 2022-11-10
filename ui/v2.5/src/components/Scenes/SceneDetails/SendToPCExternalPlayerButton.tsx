import React from "react";
import { Button } from "react-bootstrap";
import { useIntl } from "react-intl";
import { Icon } from "src/components/Shared";
import { SceneDataFragment } from "src/core/generated-graphql";
import { faPlayCircle } from "@fortawesome/free-solid-svg-icons";
import { TextUtils } from "src/utils";

export interface ISendToPCExternalPlayerButtonProps {
  scene: SceneDataFragment;
}

export const SendToPCExternalPlayerButton: React.FC<ISendToPCExternalPlayerButtonProps> = ({
  scene,
}) => {
  const { paths } = scene;
  const { external_player } = paths;
  const intl = useIntl();

  if (external_player === undefined || external_player === null || external_player === "")
  {
    return <span />
  }

  const webRequest = {
    method: 'PUT',
  };

  const clickHandler = () => {
    fetch(external_player, webRequest);
  }

  return (
    <Button
      className="minimal px-0 px-sm-2 pt-2"
      variant="secondary"
      title={intl.formatMessage({ id: "actions.open_in_external_player" })}
      onClick={clickHandler}
    >
      <Icon icon={faPlayCircle} color="white" />
    </Button>
  );
};
