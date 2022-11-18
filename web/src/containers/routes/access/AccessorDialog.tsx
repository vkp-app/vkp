import {
	Button,
	Dialog,
	DialogActions,
	DialogContent,
	DialogTitle,
	FormControlLabel,
	FormGroup,
	FormLabel,
	MenuItem,
	Select,
	Switch,
	TextField,
	Theme
} from "@mui/material";
import React, {useMemo, useState} from "react";
import {makeStyles} from "tss-react/mui";
import {AccessRef} from "../../../generated/graphql";

const useStyles = makeStyles()((theme: Theme) => ({
	formLabel: {
		fontFamily: "Manrope",
		fontWeight: "bold",
		fontSize: 13,
		paddingBottom: theme.spacing(1)
	}
}));

interface Props {
	open: boolean;
	setOpen: (v: boolean) => void;
	onAdd: (ref: AccessRef) => void;
}

const REF_USER = "User";
const REF_GROUP = "Group";
const REF_TYPES: string[] = [REF_USER, REF_GROUP];

// https://stackoverflow.com/a/8829363
const emailRegex = new RegExp(/^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/);

const AccessorDialog: React.FC<Props> = ({open, setOpen, onAdd}): JSX.Element => {
	// hooks
	const {classes} = useStyles();

	// local data
	const [name, setName] = useState<string>("");
	const [readOnly, setReadOnly] = useState<boolean>(true);
	const [refType, setRefType] = useState<string>(REF_USER);

	const emailIsValid = useMemo(() => {
		if (refType === REF_GROUP)
			return true;
		return emailRegex.test(name);
	}, [name, refType]);

	const handleAdd = (): void => {
		onAdd({
			user: refType === REF_USER ? name : "",
			group: refType === REF_GROUP ? name : "",
			readOnly: readOnly
		});
		setOpen(false);
	}

	return <Dialog
		maxWidth="sm"
		fullWidth
		open={open}
		onClose={() => setOpen(false)}
		TransitionProps={{
			onEntering: () => {
				setName(() => "");
				setReadOnly(() => true);
				setRefType(() => REF_USER);
			}
		}}>
		<DialogTitle>
			Add permission
		</DialogTitle>
		<DialogContent>
			<FormGroup>
				<FormLabel
					className={classes.formLabel}>
					{refType === REF_USER ? "Email" : "Name"}
				</FormLabel>
				<TextField
					size="small"
					value={name}
					onChange={event => setName(() => event.target.value)}
					error={refType === REF_USER ? !emailIsValid : undefined}
					helperText={refType === REF_USER ? "Must be a valid email address." : undefined}
				/>
			</FormGroup>
			<FormGroup>
				<FormLabel
					className={classes.formLabel}
					sx={{mt: 1}}>
					Type
				</FormLabel>
				<Select
					size="small"
					value={refType}
					onChange={(e) => setRefType(() => e.target.value)}>
					{REF_TYPES.map(t => <MenuItem
						key={t}
						value={t}>
						{t}
					</MenuItem>)}
				</Select>
			</FormGroup>
			<FormGroup>
				<FormControlLabel
					sx={{mt: 2}}
					labelPlacement="start"
					control={<Switch/>}
					label="Read-only"
					checked={readOnly}
					onChange={(event, checked) => setReadOnly(() => checked)}
				/>
			</FormGroup>
		</DialogContent>
		<DialogActions>
			<Button
				onClick={() => setOpen(false)}>
				Cancel
			</Button>
			<Button
				disabled={name.trim() === "" || !emailIsValid}
				onClick={handleAdd}>
				Add
			</Button>
		</DialogActions>
	</Dialog>
}
export default AccessorDialog;
