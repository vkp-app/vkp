/*
 *    Copyright 2022 Django Cass
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 *
 */

import React, {useState} from "react";
import {
	AppBar,
	Avatar,
	ButtonBase,
	Divider,
	IconButton,
	ListItem,
	ListItemButton,
	ListItemText,
	MenuItem,
	Popover,
	Theme,
	Toolbar,
	Typography,
} from "@mui/material";
import {Link} from "react-router-dom";
import {useTheme} from "@mui/material/styles";
import {makeStyles} from "tss-react/mui";
import {ChevronDown, Help, User} from "tabler-icons-react";
import {useCurrentUserQuery} from "../generated/graphql";

const useStyles = makeStyles()((theme: Theme) => ({
	grow: {
		flexGrow: 1
	},
	brandButton: {
		borderRadius: theme.spacing(1),
		paddingRight: theme.spacing(1),
		height: 40
	},
	brand: {
		paddingRight: 8,
		[theme.breakpoints.up("sm")]: {
			paddingRight: 0
		},
		fontFamily: "Manrope",
		fontWeight: 500,
		pointerEvents: "none"
	},
	title: {
		display: "none",
		[theme.breakpoints.up("sm")]: {
			display: "block"
		},
		fontFamily: "Manrope",
		pointerEvents: "none"
	},
	sectionDesktop: {
		marginRight: theme.spacing(1),
		display: "none",
		[theme.breakpoints.up("md")]: {
			display: "flex"
		},
		alignItems: "center"
	},
	menuIcon: {
		paddingRight: theme.spacing(1)
	},
	avatar: {
		width: 24,
		height: 24,
		margin: theme.spacing(1.5),
		borderRadius: 0
	},
	menuButton: {
		marginRight: 36,
		marginLeft: theme.spacing(0.5)
	},
	hide: {
		display: "none",
	},
	toolbar: {
		display: "flex",
		alignItems: "center",
		justifyContent: "flex-end",
		padding: theme.spacing(0, 1)
	},
}));

interface NavProps {
	loading?: boolean;
}

const Nav: React.FC<NavProps> = ({loading = false}: NavProps): JSX.Element => {
	// hooks
	const {classes} = useStyles();
	const theme = useTheme();
	const {data} = useCurrentUserQuery();

	// local state
	const [anchorEl, setAnchorEl] = useState<HTMLElement | null>(null);

	const handleMenuClose = (): void => {
		setAnchorEl(null);
	};

	return (
		<div>
			<AppBar
				elevation={0}
				variant="outlined"
				position="fixed"
				color="inherit">
				<Toolbar
					className={classes.toolbar}
					variant="dense">
					<ButtonBase
						className={classes.brandButton}
						component={Link}
						to="/">
						<Avatar
							className={classes.avatar}
							src="/src/img/kubernetes-icon-color.svg"
							alt="VKP logo"
						/>
						<Typography
							className={classes.brand}
							variant="h6"
							color="textPrimary">
							VKP
						</Typography>
						<Typography
							className={classes.brand}
							sx={{ml: 1}}
							variant="h6"
							color="textSecondary">
							Virtual Kubernetes Platform
						</Typography>
					</ButtonBase>
					<div className={classes.grow}/>
					<div className={classes.sectionDesktop}>
						<IconButton
							style={{margin: 8}}
							disabled={loading}
							component={Link}
							centerRipple={false}
							size="small"
							color="inherit"
							to="/help/overview">
							<Help color={theme.palette.text.secondary}/>
						</IconButton>
						<ButtonBase
							className={classes.brandButton}
							sx={{pl: 1, pr: 1}}
							focusRipple
							disabled={loading}
							onClick={e => setAnchorEl(e.currentTarget)}>
							<User
								size={22}
								color={theme.palette.text.secondary}
							/>
							<ChevronDown
								style={{marginLeft: theme.spacing(1)}}
								size={16}
								color={theme.palette.text.secondary}
							/>
						</ButtonBase>
					</div>
				</Toolbar>
			</AppBar>
			<Popover
				sx={{mt: 1.5}}
				PaperProps={{variant: "outlined", elevation: 0, sx: {minWidth: 200}}}
				anchorEl={anchorEl}
				anchorOrigin={{vertical: "bottom", horizontal: "right"}}
				transformOrigin={{vertical: "top", horizontal: "right"}}
				open={anchorEl != null && !loading}
				onClose={handleMenuClose}>
				{!data && <ListItemButton
					component="a"
					href="/auth/redirect"
					rel="noopener noreferrer">
					<ListItemText
						primary="Not logged in."
						secondary="Click here to login"
					/>
				</ListItemButton>}
				{data && <ListItem>
					<ListItemText primary={data?.currentUser.username || ""}/>
				</ListItem>}
				<Divider/>
				<MenuItem
					sx={{fontSize: 14}}
					onClick={handleMenuClose}
					component={Link}
					to="/profile"
					disabled={!data}>
					Edit profile
				</MenuItem>
				<MenuItem
					sx={{fontSize: 14}}
					onClick={handleMenuClose}
					component={Link}
					to="/profile/preferences"
					disabled={!data}>
					Preferences
				</MenuItem>
			</Popover>
		</div>
	);
};
export default Nav;
