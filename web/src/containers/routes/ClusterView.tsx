import React from "react";
import StandardLayout from "../layout/StandardLayout";
import {Link, useParams} from "react-router-dom";
import {useClusterQuery} from "../../generated/graphql";
import {
	Button,
	Card,
	CardHeader, IconButton,
	List,
	ListItem, ListItemSecondaryAction,
	ListItemText,
	ListSubheader,
	Skeleton,
	Theme,
	Typography
} from "@mui/material";
import InlineError from "../alert/InlineError";
import {makeStyles} from "tss-react/mui";
import {ArrowLeft, ExternalLink} from "tabler-icons-react";

const useStyles = makeStyles()((theme: Theme) => ({
	title: {
		fontFamily: "Manrope",
		fontSize: 24,
		fontWeight: 500,
		paddingBottom: theme.spacing(1)
	},
	button: {
		fontFamily: "Manrope",
		fontWeight: 600,
		fontSize: 12,
		textTransform: "none",
		minHeight: 24,
		height: 24
	}
}));

const ClusterView: React.FC = (): JSX.Element => {
	// hooks
	const params = useParams();
	const {classes} = useStyles();

	const clusterName = params["name"];
	const tenantName = params["tenant"];

	const {data, loading, error} = useClusterQuery({
		variables: {tenant: tenantName || "", cluster: clusterName || ""},
		skip: !clusterName || !tenantName
	});

	return <StandardLayout>
		<ListSubheader
			sx={{display: "flex", alignItem: "center"}}>
			<IconButton
				size={"small"}
				centerRipple={false}
				component={Link}
				to={`/clusters/${tenantName}`}>
				<ArrowLeft
					size={18}
				/>
			</IconButton>
			Cluster dashboard
		</ListSubheader>
		<Card
			variant={"outlined"}
			sx={{p: 2}}>
			{!loading && error && <InlineError
				message={"Unable to load cluster information"}
				error={error}
			/>}
			<List>
				<ListItem>
					<ListItemText
						primaryTypographyProps={{className: classes.title}}
						primary={loading ? <Skeleton variant={"text"} height={30} width={"40%"}/> : data?.cluster.name}
						secondary={loading ? <Skeleton variant={"text"} height={20} width={"20%"}/> : `Kubernetes v${data?.cluster.status.kubeVersion || "1.XX"}`}
					/>
				</ListItem>
				<ListItem>
					<ListItemText
						primary={"Kubernetes API address"}
						secondary={loading ? <Skeleton variant={"text"} height={20} width={"40%"}/> : `https://${data?.cluster.status.kubeURL}:443`}
					/>
					<ListItemSecondaryAction>
						<Button
							className={classes.button}
							disabled={loading}
							variant={"outlined"}
							startIcon={<ExternalLink
								size={18}
							/>}
							href={`https://${data?.cluster.status.kubeURL}`}
							target={"_blank"}>
							Open
						</Button>
					</ListItemSecondaryAction>
				</ListItem>
				<ListItem>
					<ListItemText
						primary={"Management API"}
						secondary={loading ? <Skeleton variant={"text"} height={20} width={"40%"}/> : "VCluster"}
					/>
				</ListItem>
			</List>
		</Card>
	</StandardLayout>
}
export default ClusterView;
