import {
	Box,
	Button,
	Card,
	CardHeader,
	Checkbox,
	FormControl,
	FormControlLabel,
	FormGroup,
	FormLabel,
	Switch
} from "@mui/material";
import React, {useEffect, useMemo, useState} from "react";
import {useParams} from "react-router-dom";
import cronstrue from "cronstrue";
import StandardLayout from "../../layout/StandardLayout";
import BackButton from "../../layout/BackButton";
import {useMaintenancePolicyQuery, useSetMaintenancePolicyMutation} from "../../../generated/graphql";
import InlineError from "../../alert/InlineError";

const DAYS = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"];

const ClusterMaintenanceWindow: React.FC = (): JSX.Element => {
	// hooks
	const params = useParams();

	const clusterName = params["name"] || "";
	const tenantName = params["tenant"] || "";

	// global state
	const maintenancePolicy = useMaintenancePolicyQuery({
		variables: {tenant: tenantName, cluster: clusterName},
		skip: !tenantName || !clusterName
	});
	const [setMaintenancePolicy, setMaintenancePolicyData] = useSetMaintenancePolicyMutation();

	// local state
	const [windowEnabled, setWindowEnabled] = useState<boolean>(false);
	const [selectedDays, setSelectedDays] = useState<string[]>([]);

	useEffect(() => {
		if (!maintenancePolicy.data) {
			setWindowEnabled(() => false);
			setSelectedDays(() => []);
			return;
		}
		const cronSegments = maintenancePolicy.data.clusterMaintenanceWindow.schedule.split(" ");
		if (cronSegments[cronSegments.length - 1] === "*") {
			setWindowEnabled(() => false);
			setSelectedDays(() => []);
			return;
		}
		const daysOfWeek = cronSegments[cronSegments.length - 1].split(",").filter(s => s !== "");
		setSelectedDays(() => daysOfWeek);
		setWindowEnabled(() => true);
	}, [maintenancePolicy.data]);

	const dayShortCode = (s: string): string => {
		return s.toLocaleUpperCase().substring(0, 3);
	}

	const cron = useMemo(() => {
		if (!windowEnabled || selectedDays.length === 0)
			return "* * * * *";
		return `* * * * ${selectedDays.join(",")}`;
	}, [windowEnabled, selectedDays]);

	const cronExpl = cronstrue.toString(cron);

	const handleSave = (): void => {
		void setMaintenancePolicy({
			variables: {tenant: tenantName, cluster: clusterName, schedule: cron}
		});
	}
	
	const loading = maintenancePolicy.loading || setMaintenancePolicyData.loading;
	
	return <StandardLayout>
		<BackButton
			title="Back to cluster"
			to={`/clusters/${tenantName}/cluster/${clusterName}`}
		/>
		<Card
			sx={{p: 2}}>
			<CardHeader
				title="Maintenance policy"
				titleTypographyProps={{fontFamily: "Figtree", fontWeight: 400, mb: 1}}
				subheader="Control when upgrades and maintenance can take place. Normally VKP will perform upgrades at any time."
				subheaderTypographyProps={{fontSize: 14}}
			/>
		</Card>
		<Card
			sx={{p: 2, mt: 2}}>
			<CardHeader
				title="Enable maintenance window"
				titleTypographyProps={{variant: "body1"}}
				subheader={cron !== "* * * * *" ? `Upgrades may be performed ${cronExpl.charAt(0).toLocaleLowerCase()}${cronExpl.slice(1)}.` : "Upgrades may be performed at any time."}
				avatar={<Switch
					checked={windowEnabled}
					onChange={(e, checked) => setWindowEnabled(() => checked)}
					disabled={loading}
				/>}
			/>
			<FormControl
				sx={{m: 2}}>
				<FormLabel
					component="legend">
					Days
				</FormLabel>
				<FormGroup
					row>
					{DAYS.map(d => {
						const sc = dayShortCode(d);
						return <FormControlLabel
							key={d}
							control={<Checkbox/>}
							label={d}
							disabled={loading}
							checked={selectedDays.includes(sc)}
							value={selectedDays.includes(sc)}
							onChange={(e, checked) => checked ? setSelectedDays(s => [...s, sc]) : setSelectedDays(s => s.filter(v => v !== sc))}
						/>
					})}
				</FormGroup>
			</FormControl>
			<InlineError
				error={maintenancePolicy.error || setMaintenancePolicyData.error}
			/>
		</Card>
		<Box
			sx={{display: "flex", p: 1, pt: 2}}>
			<Box sx={{flexGrow: 1}}/>
			<Button
				sx={{ml: 1}}
				disabled={loading}
				onClick={handleSave}>
				Update
			</Button>
		</Box>
	</StandardLayout>
}
export default ClusterMaintenanceWindow;
