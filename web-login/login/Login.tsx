import {Avatar, Box, Button, Card, List, ListItem, ListItemText, Theme, Typography} from "@mui/material";
import React from "react";
import {default as Grid} from "@mui/material/Unstable_Grid2/Grid2";
import {makeStyles} from "tss-react/mui";

const useStyles = makeStyles()((theme: Theme) => ({
	container: {
		height: "100vh",
		display: "flex",
		alignItems: "center",
		backgroundColor: theme.palette.background.default
	},
	item: {
		display: "flex"
	},
	button: {
		fontFamily: "Figtree",
		textTransform: "none",
		textAlign: "center",
		marginBottom: theme.spacing(1)
	},
	avatar: {
		width: 56,
		height: 56,
		margin: theme.spacing(1.5),
		borderRadius: 0
	},
	brand: {
		fontFamily: "Figtree",
		fontWeight: 500,
	},
}));

const Login: React.FC = (): JSX.Element => {
	// hooks
	const {classes} = useStyles();

	return <div>
		<Grid
			className={classes.container}
			container
			spacing={2}>
			<Grid xs={6}>
				<Grid container>
					<Grid xs={6}/>
					<Grid xs={6}>
						<Card
							sx={{m: 1, p: 1}}>
							<Box
								sx={{p: 4}}>
								<Typography
									sx={{fontFamily: "Figtree"}}
									variant="h5">
									Login with...
								</Typography>
								<List
									sx={{mt: 1}}>
									<ListItem
										className={classes.button}
										component={Button}>
										<ListItemText>
											GitLab
										</ListItemText>
									</ListItem>
									<ListItem
										className={classes.button}
										component={Button}>
										<ListItemText>
											GitHub
										</ListItemText>
									</ListItem>
								</List>
							</Box>
						</Card>
					</Grid>
				</Grid>
			</Grid>
			<Grid
				className={classes.item}
				xs={6}>
				<Avatar
					className={classes.avatar}
					src={new URL("/kubernetes-icon-color.svg", import.meta.url).href}
					alt="VKP logo"
				/>
				<Box>
					<Typography
						className={classes.brand}
						variant="h4"
						color="primary.main">
						VKP
					</Typography>
					<Typography
						className={classes.brand}
						variant="h6"
						color="text.secondary">
						Virtual Kubernetes Platform
					</Typography>
				</Box>
			</Grid>
		</Grid>
	</div>
}
export default Login;
