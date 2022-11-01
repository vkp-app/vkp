import React from 'react'
import './App.css'
import createCache from "@emotion/cache";
import {CacheProvider} from "@emotion/react";
import {makeStyles} from "tss-react/mui";
import {CssBaseline, Theme} from "@mui/material";
import Nav from "./containers/Nav";

const useStyles = makeStyles()((theme: Theme) => ({
	root: {
		display: "flex",
	},
	toolbar: {
		display: "flex",
		alignItems: "center",
		justifyContent: "flex-end",
		padding: theme.spacing(0, 1),
		height: 48
	},
	content: {
		flexGrow: 1
	},
}));

const App: React.FC = (): JSX.Element => {
	// hooks
	const {classes} = useStyles();

	const cache = createCache({
		key: "mui",
		prepend: true
	});

	return (
		<CacheProvider value={cache}>
			<div className={classes.root}>
				<CssBaseline/>
				<Nav/>
			</div>
		</CacheProvider>
	)
}

export default App;
