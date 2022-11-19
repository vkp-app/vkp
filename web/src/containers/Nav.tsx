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
	Fade,
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
import Icon from "@mdi/react";
import {mdiAccountCircle, mdiChevronDown, mdiHelpCircle} from "@mdi/js";
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
		fontFamily: "Figtree",
		fontWeight: 500,
		pointerEvents: "none"
	},
	title: {
		display: "none",
		[theme.breakpoints.up("sm")]: {
			display: "block"
		},
		fontFamily: "Figtree",
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
		width: 32,
		height: 32,
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
	const [expanded, setExpanded] = useState<boolean>(false);

	const handleMenuClose = (): void => {
		setAnchorEl(null);
	};

	return (
		<div>
			<AppBar
				elevation={2}
				position="fixed"
				color="inherit">
				<Toolbar
					className={classes.toolbar}
					variant="dense">
					<ButtonBase
						onMouseEnter={() => setExpanded(() => true)}
						onMouseLeave={() => setExpanded(() => false)}
						className={classes.brandButton}
						component={Link}
						to="/tenants">
						<Avatar
							className={classes.avatar}
							src="/src/img/kubernetes-icon-color.svg"
							alt="VKP logo"
						/>
						<Typography
							className={classes.brand}
							variant="h6"
							color="primary.main">
							VKP
						</Typography>
						<Fade
							in={expanded}>
							<Typography
								className={classes.brand}
								sx={{ml: 1}}
								variant="h6"
								color="text.secondary">
								Virtual Kubernetes Platform
							</Typography>
						</Fade>
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
							<Icon
								path={mdiHelpCircle}
								size={1}
								color={theme.palette.primary.main}
							/>
						</IconButton>
						<ButtonBase
							className={classes.brandButton}
							sx={{pl: 1, pr: 1}}
							focusRipple
							disabled={loading}
							onClick={e => setAnchorEl(e.currentTarget)}>
							<Icon
								path={mdiAccountCircle}
								size={1}
								color={theme.palette.primary.main}
							/>
							<Icon
								style={{marginLeft: theme.spacing(0.5)}}
								path={mdiChevronDown}
								size={1}
								color={theme.palette.primary.main}
							/>
						</ButtonBase>
					</div>
				</Toolbar>
			</AppBar>
			<Popover
				sx={{mt: 1.5}}
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
				<MenuItem
					sx={{fontSize: 14}}
					onClick={handleMenuClose}
					component={Link}
					to="/about">
					About
				</MenuItem>
			</Popover>
		</div>
	);
};
export default Nav;
