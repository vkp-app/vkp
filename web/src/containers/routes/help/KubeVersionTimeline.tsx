import React, {useMemo} from "react";
import {
	Timeline,
	TimelineConnector,
	TimelineContent,
	TimelineDot,
	TimelineItem,
	TimelineOppositeContent,
	TimelineSeparator
} from "@mui/lab";
import {Button, Card, CardHeader, Collapse, Typography} from "@mui/material";
import {useLocation, useNavigate} from "react-router-dom";
import Icon from "@mdi/react";
import {mdiOpenInNew} from "@mdi/js";
import StandardLayout from "../../layout/StandardLayout";

interface KubeVersion {
	version: number;
	name: string;
	description: string;
	date: string;
	link?: string;
}

const KUBE_VERSIONS: KubeVersion[] = [
	{
		version: 27,
		name: "Chill Vibes",
		description: "In-place resource updates, VolumeGroupSnapshot, registry.k8s.io",
		date: "Q2 2023",
		link: "https://kubernetes.io/blog/2023/04/11/kubernetes-v1-27-release/"
	},
	{
		version: 26,
		name: "Electrifying",
		description: "Control-plane metrics, Pod scheduling improvements.",
		date: "Q4 2022",
		link: "https://kubernetes.io/blog/2022/12/09/kubernetes-v1-26-release/",
	},
	{
		version: 25,
		name: "Combiner",
		description: "PodSecurity, cgroups v2, CSIMigration, KMSv2.",
		date: "Q3 2022",
		link: "https://kubernetes.io/blog/2022/08/23/kubernetes-v1-25-release/"
	},
	{
		version: 24,
		name: "Stargazer",
		description: "End of an Era - removal of the DockerShim.",
		date: "Q2 2022",
		link: "https://kubernetes.io/blog/2022/05/03/kubernetes-1-24-release-announcement/"
	},
	{
		version: 23,
		name: "The Next Frontier",
		description: "IPv4/IPv6 Dual Stack, HPAv2, PodSecurity.",
		date: "Q4 2021",
		link: "https://kubernetes.io/blog/2021/12/07/kubernetes-1-23-release-announcement/"
	},
	{
		version: 22,
		name: "Reaching New Peaks",
		description: "Server-side Apply, External credential providers",
		date: "Q3 2021",
		link: "https://kubernetes.io/blog/2021/08/04/kubernetes-1-22-release-announcement/"
	},
	{
		version: 21,
		name: "Power to the Community",
		description: "CronJobs, Immutable Secrets, PV Health",
		date: "Q2 2021",
		link: "https://kubernetes.io/blog/2021/04/08/kubernetes-1-21-release-announcement/"
	}
];

const KubeVersionTimeline: React.FC = (): JSX.Element => {
	// hooks
	const navigate = useNavigate();
	const location = useLocation();

	// local state
	const selected = useMemo(() => {
		if (!location.hash || location.hash === "#")
			return 0;
		return Number(location.hash.slice(1));
	}, [location.hash]);

	const handleSelect = (n: number): void => {
		navigate({
			...location,
			hash: `#${n}`
		});
	}

	return <StandardLayout>
		<Card
			sx={{p: 2}}>
			<CardHeader
				title="Kubernetes version timeline"
				titleTypographyProps={{fontFamily: "Figtree", fontWeight: 400, mb: 1}}
				subheader={<span>
					Kubernetes releases 3 times per year which can be difficult to keep up with.
					This page provides a general overview of which Kubernetes versions are available and the headlining features that came with them.
				</span>}
				subheaderTypographyProps={{fontSize: 14}}
			/>
		</Card>
		<Timeline
			position="alternate">
			{KUBE_VERSIONS.map(v => {
				const isSelected = selected === v.version;
				const notUnselected = selected > 0 && !isSelected;
				return <TimelineItem
					key={v.name}>
					<TimelineOppositeContent
						onClick={() => handleSelect(v.version)}>
						<Typography
							color={notUnselected ? "text.secondary" : undefined}
							sx={{fontFamily: "monospace, monospace"}}>
							v1.{v.version}
						</Typography>
						<Typography
							color="text.secondary"
							variant="caption">
							{v.date}
						</Typography>
					</TimelineOppositeContent>
					<TimelineSeparator>
						<TimelineConnector/>
						<TimelineDot
							variant={isSelected ? "filled" : "outlined"}
							onClick={() => handleSelect(v.version)}
							color={isSelected ? "secondary" : "primary"}
						/>
					</TimelineSeparator>
					<TimelineContent
						onClick={() => handleSelect(v.version)}>
						<Typography
							variant="h6"
							color={notUnselected ? "text.secondary" : undefined}
							sx={{fontFamily: "Figtree", fontWeight: isSelected ? 600 : 500}}>
							{v.name}
						</Typography>
						<Typography
							color={notUnselected ? "text.secondary" : undefined}>
							{v.description}
						</Typography>
						{v.link && <Collapse
							in={isSelected}>
							<Button
								variant="text"
								href={v.link}
								target="_blank"
								endIcon={<Icon
									path={mdiOpenInNew}
									size={0.7}
								/>}>
								More information
							</Button>
						</Collapse>}
					</TimelineContent>
				</TimelineItem>
			})}
		</Timeline>
	</StandardLayout>
}
export default KubeVersionTimeline;
