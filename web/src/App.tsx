import React from 'react'
import './App.css'
import createCache from "@emotion/cache";
import {CacheProvider} from "@emotion/react";
import {makeStyles} from "tss-react/mui";
import {CssBaseline, Theme} from "@mui/material";
import Nav from "./containers/Nav";
import {Route, Routes} from "react-router-dom";
import ClusterList from "./containers/routes/ClusterList";
import NotFound from "./containers/routes/NotFound";
import ClusterView from "./containers/routes/ClusterView";
import Home from "./containers/routes/Home";
import CreateCluster from "./containers/routes/CreateCluster";

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
				<main className={classes.content}>
					<div className={classes.toolbar}/>
					<Routes>
						<Route
							path={"/"}
							element={<Home/>}
						/>
						<Route
							path={"/clusters/:tenant/cluster/:name"}
							element={<ClusterView/>}
						/>
						<Route
							path={"/clusters/:tenant"}
							element={<ClusterList/>}
						/>
						<Route
							path={"/new/cluster"}
							element={<CreateCluster/>}
						/>
						<Route
							path={"/*"}
							element={<NotFound/>}
						/>
					</Routes>
				</main>
			</div>
		</CacheProvider>
	)
}

export default App;
