import {Box, IconButton, Link, Tooltip} from "@mui/material";
import React, {ReactNode, useEffect, useMemo, useState} from "react";
import {AlertCircle, CircleCheck, InfoCircle} from "tabler-icons-react";
import {useTheme} from "@mui/material/styles";
import {KUBERNETES_VERSION_LATEST, KUBERNETES_VERSION_MAX, KUBERNETES_VERSION_MIN} from "../../../config/constants";

interface Props {
	version: string | number;
	showLabel?: boolean;
}

const ClusterVersionIndicator: React.FC<Props> = ({version, showLabel = false}): JSX.Element => {
	// hooks
	const theme = useTheme();

	// local state
	const [tooltip, setTooltip] = useState<string>("");
	const [icon, setIcon] = useState<ReactNode>(<div></div>);

	const num = useMemo(() => {
		if (typeof version === "string") {
			if (version.startsWith("1.")) {
				return Number(version.replace("1.", ""));
			}
			return Number(version);
		}
		return version;
	}, [version]);

	useEffect(() => {
		if (num === KUBERNETES_VERSION_LATEST) {
			setTooltip(() => "Cluster is up-to-date");
			setIcon(() => <CircleCheck size={18} color={theme.palette.success.main}/>);
		}
		else if (num > KUBERNETES_VERSION_LATEST && num <= KUBERNETES_VERSION_MAX) {
			setTooltip(() => "Cluster is ahead of the stable release");
			setIcon(() => <InfoCircle size={18} color={theme.palette.info.main}/>);
		}
		else if (num >= KUBERNETES_VERSION_MIN && num < KUBERNETES_VERSION_LATEST) {
			setTooltip(() => "Cluster is behind the stable release");
			setIcon(() => <InfoCircle size={18} color={theme.palette.warning.main}/>);
		}
		else if (num > 0 && num < KUBERNETES_VERSION_MIN) {
			setTooltip(() => "Cluster is running an un-supported version of Kubernetes");
			setIcon(() => <AlertCircle size={18} color={theme.palette.error.main}/>);
		} else {
			setTooltip(() => "Unknown or unsupported version");
			setIcon(() => <AlertCircle size={18} color={theme.palette.text.secondary}/>);
		}
	}, [num]);

	return <Box
		sx={{display: "flex", alignItems: "center"}}>
		{showLabel && <span>Kubernetes&nbsp;</span>}
		<Link
			href={`https://kubernetes.io/releases/#release-v1-${num}`}
			target="_blank">
			v1.{num}
		</Link>
		<Tooltip title={tooltip}>
			<IconButton
				sx={{ml: 1}}
				size="small"
				centerRipple={false}>
				{icon}
			</IconButton>
		</Tooltip>
	</Box>
}
export default ClusterVersionIndicator;
