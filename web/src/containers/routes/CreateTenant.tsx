import React, {useMemo, useState} from "react";
import {Box, Button, Card, CardHeader, Divider, FormGroup, FormLabel, TextField, Theme} from "@mui/material";
import {makeStyles} from "tss-react/mui";
import {useNavigate} from "react-router-dom";
import StandardLayout from "../layout/StandardLayout";
import {useCreateTenantMutation} from "../../generated/graphql";
import InlineError from "../alert/InlineError";

const useStyles = makeStyles()((theme: Theme) => ({
	button: {
		marginRight: theme.spacing(1)
	},
	formLabel: {
		fontFamily: "Manrope",
		fontWeight: "bold",
		fontSize: 13,
		paddingBottom: theme.spacing(1)
	}
}));

const nameRegex = new RegExp(/^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/);

const CreateTenant: React.FC = (): JSX.Element => {
	// hooks
	const {classes} = useStyles();
	const navigate = useNavigate();

	const [createTenant, {loading, error}] = useCreateTenantMutation();

	// local state
	const [name, setName] = useState<string>("");

	const nameIsValid = useMemo(() => {
		return nameRegex.test(name);
	}, [name]);

	const onCreate = (): void => {
		createTenant({
			variables: {tenant: name},
		}).then(r => {
			if (!r.errors) {
				navigate(`/clusters/${name}`);
			}
		});
	}

	return <StandardLayout>
		<Card
			sx={{p: 2, pl: 0, pt: 0}}
			variant="outlined">
			<CardHeader
				title="Request a tenancy"
			/>
			{!loading && error && <InlineError
				message="Unable to lodge tenant creation request"
				error={error}
			/>}
			<Divider/>
			<Box
				sx={{ml: 4, mr: 2, mt: 2}}>
				<FormGroup>
					<FormLabel
						className={classes.formLabel}>
								Name
					</FormLabel>
					<TextField
						size="small"
						value={name}
						onChange={e => setName(() => e.target.value)}
						error={!nameIsValid}
						helperText="Must consist of lower case alphanumeric characters or '-', and must start and end with an alphanumeric character."
					/>
				</FormGroup>
				<Box
					sx={{pt: 2, display: "flex", float: "right"}}>
					<Button
						className={classes.button}
						disabled={loading}
						onClick={() => navigate(-1)}>
							Cancel
					</Button>
					<Button
						className={classes.button}
						disabled={!nameIsValid || loading}
						onClick={onCreate}>
						Create
					</Button>
				</Box>
			</Box>
		</Card>
	</StandardLayout>
}
export default CreateTenant;
