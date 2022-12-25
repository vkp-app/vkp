import React from "react";
import {
	Avatar,
	Button,
	Card,
	CardActions,
	CardContent,
	CardHeader,
	CircularProgress,
	Skeleton,
	Typography
} from "@mui/material";
import Grid2 from "@mui/material/Unstable_Grid2";
import {AddonPhase, AddonSource, ClusterAddon} from "../../../generated/graphql";
import AddonSourceChip from "./AddonSourceChip";
import AddonChip from "./AddonChip";


interface Props {
	item: ClusterAddon | null;
	phase?: AddonPhase | null;
	loading?: boolean;
	onInstall: () => void;
	onUninstall: () => void;
	readOnly?: boolean;
}

const AddonItem: React.FC<Props> = ({item, phase, loading, readOnly, onInstall, onUninstall}): JSX.Element => {
	const handleClick = (): void => {
		if (phase != null) {
			onUninstall();
			return;
		}
		onInstall();
	}

	const getAction = (): string => {
		switch (phase) {
			case undefined:
			case null:
				return "Install";
			case AddonPhase.Installed:
				return "Uninstall";
			default:
				return phase.toString();
		}
	}

	return <Grid2
		lg={6}
		md={12}>
		<Card
			sx={{maxHeight: 250, minHeight: 250, p: 1}}>
			<CardHeader
				disableTypography={item == null}
				title={item?.displayName || <Skeleton
					variant="text"
					width="40%"
				/>}
				subheader={item?.maintainer ? <span>
					{item.maintainer}
					<AddonSourceChip source={item?.source || AddonSource.Unknown}/>
					{phase != null && <AddonChip/>}
				</span> : <Skeleton
					variant="text"
					width="60%"
				/>}
				titleTypographyProps={{fontFamily: "Figtree", fontWeight: "500", fontSize: 16}}
				subheaderTypographyProps={{fontSize: 14}}
				avatar={item != null ? <Avatar
					src={item.logo}
					alt={`${item.displayName} logo`}
					variant={item.logo !== "" ? "square" : "circular"}
				/> : <Skeleton
					variant="circular"
					width={48}
					height={48}
				/>}
				action={item != null ? <Button
					variant="outlined"
					disabled={loading || phase === AddonPhase.Installing || phase === AddonPhase.Deleting || readOnly}
					startIcon={loading || phase === AddonPhase.Installing || phase === AddonPhase.Deleting ? <CircularProgress size={14}/> : undefined}
					onClick={handleClick}>
					{getAction()}
				</Button> : undefined}
			/>
			<CardContent>
				{item != null ? <Typography
					variant="body2"
					color="text.secondary">
					{item.description || <i>No description provided.</i>}
				</Typography> : <Skeleton
					variant="text"
					width="40%"
				/>}
			</CardContent>
			{(item == null || item?.sourceURL !== "") && <CardActions>
				{item != null ? <Button
					variant="text"
					href={item?.sourceURL || ""}
					target="_blank">
					View source
				</Button> : <Skeleton
					width="20%"
				/>}
			</CardActions>}
		</Card>
	</Grid2
>
}
export default AddonItem;
