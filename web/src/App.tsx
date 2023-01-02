import React, {useEffect, useMemo} from "react"
import createCache from "@emotion/cache";
import {CacheProvider} from "@emotion/react";
import {makeStyles} from "tss-react/mui";
import {Route, Routes} from "react-router-dom";
import {CssBaseline, Theme, useMediaQuery} from "@mui/material";
import {createTheme, ThemeProvider} from "@mui/material/styles";
import Nav from "./containers/Nav";
import ClusterList from "./containers/routes/ClusterList";
import NotFound from "./containers/routes/NotFound";
import ClusterView from "./containers/routes/ClusterView";
import Home from "./containers/routes/Home";
import CreateCluster from "./containers/routes/CreateCluster";
import TenantList from "./containers/routes/TenantList";
import AddonList from "./containers/routes/AddonList";
import CreateTenant from "./containers/routes/CreateTenant";
import ClusterAccessorList from "./containers/routes/ClusterAccessorList";
import TenantAccessorList from "./containers/routes/TenantAccessorList";
import About from "./containers/routes/About";
import {FF_BANNER_ENABLED, FF_BANNER_HEIGHT} from "./config/constants";
import ClusterMaintenanceWindow from "./containers/routes/cluster/ClusterMaintenanceWindow";

const useStyles = makeStyles()((theme: Theme) => ({
	root: {
		display: "flex",
	},
	toolbar: {
		display: "flex",
		alignItems: "center",
		justifyContent: "flex-end",
		padding: theme.spacing(0, 1),
		height: 48 + (FF_BANNER_ENABLED ? FF_BANNER_HEIGHT : 0)
	},
	content: {
		flexGrow: 1
	},
}));

export const DEFAULT_LOAD_DELAY_MS = 200;
// todo convert to feature flags (see #28)
export const FF_USE_OUTLINED_MODE = false;
export const FF_METRIC_BASE_ZERO = false;

const App: React.FC = (): JSX.Element => {
	// hooks
	const {classes} = useStyles();

	const cache = createCache({
		key: "mui",
		prepend: true
	});

	const prefersDarkMode = useMediaQuery("(prefers-color-scheme: dark)");

	const theme = useMemo(() => {
		return createTheme({
			palette: {
				mode: prefersDarkMode ? "dark" : "light",
				background: {
					default: prefersDarkMode ? "#1d1d1d" : "#F0F0F4",
					paper: prefersDarkMode ? "#2d2d2d" : "#ffffff"
				}
			},
			components: {
				MuiCheckbox: {
					defaultProps: {
						centerRipple: false
					}
				},
				MuiIconButton: {
					defaultProps: {
						centerRipple: false
					}
				},
				MuiAppBar: {
					defaultProps: {
						variant: FF_USE_OUTLINED_MODE ? "outlined" : "elevation"
					}
				},
				MuiCard: {
					defaultProps: {
						variant: FF_USE_OUTLINED_MODE ? "outlined" : "elevation"
					},
					styleOverrides: {
						root: ({theme, ownerState}) => ({
							borderRadius: ownerState.variant === "outlined" ? theme.spacing(2) : undefined
						})
					}
				},
				MuiSkeleton: {
					defaultProps: {
						animation: "wave"
					}
				},
				MuiButton: {
					defaultProps: {
						variant: "outlined"
					},
					styleOverrides: {
						root: ({theme, ownerState}) => ({
							fontFamily: "Figtree",
							fontWeight: 500,
							fontSize: 14,
							textTransform: "none",
							minHeight: 30,
							height: 30,
							borderRadius: ownerState.variant === "outlined" ? theme.spacing(1) : undefined
						})
					}
				},
				MuiDialogTitle: {
					styleOverrides: {
						root: {
							fontFamily: "Figtree",
							fontWeight: 500
						}
					}
				},
				MuiTooltip: {
					styleOverrides: {
						tooltip: ({theme}) => ({
							backgroundColor: theme.palette.background.paper,
							color: theme.palette.text.primary,
							boxShadow: theme.shadows[1],
							fontSize: "0.9rem",
							fontWeight: 400
						})
					}
				},
				MuiListSubheader: {
					styleOverrides: {
						root: {
							backgroundColor: "transparent"
						}
					}
				},
				MuiOutlinedInput: {
					styleOverrides: {
						root: {
							borderRadius: 8
						},
					}
				},
				MuiInputLabel: {
					styleOverrides: {
						shrink: {
							color: "text.primary",
							fontWeight: 500
						}
					}
				},
				MuiPopover: {
					defaultProps: {
						PaperProps: {
							variant: FF_USE_OUTLINED_MODE ? "outlined" : "elevation",
							elevation: FF_USE_OUTLINED_MODE ? 0 : 4,
							sx: {minWidth: 200}
						}
					}
				}
			}
		});
	}, [prefersDarkMode]);

	useEffect(() => {
		document.documentElement.setAttribute("data-theme", theme.palette.mode);
	}, [theme.palette.mode]);

	return (
		<CacheProvider value={cache}>
			<ThemeProvider theme={theme}>
				<div className={classes.root}>
					<CssBaseline/>
					<Nav/>
					<main className={classes.content}>
						<div className={classes.toolbar}/>
						<Routes>
							<Route
								path="/"
								element={<Home/>}
							/>
							<Route
								path="/tenants"
								element={<TenantList/>}
							/>
							<Route
								path="/clusters/:tenant/cluster/:name/-/addons"
								element={<AddonList/>}
							/>
							<Route
								path="/clusters/:tenant/cluster/:name/-/accessors"
								element={<ClusterAccessorList/>}
							/>
							<Route
								path="/clusters/:tenant/cluster/:name/-/maintenance"
								element={<ClusterMaintenanceWindow/>}
							/>
							<Route
								path="/clusters/:tenant/cluster/:name"
								element={<ClusterView/>}
							/>
							<Route
								path="/clusters/:tenant"
								element={<ClusterList/>}
							/>
							<Route
								path="/clusters/:tenant/-/accessors"
								element={<TenantAccessorList/>}
							/>
							<Route
								path="/new/cluster/:tenant"
								element={<CreateCluster/>}
							/>
							<Route
								path="/new/tenant"
								element={<CreateTenant/>}
							/>
							<Route
								path="/about"
								element={<About/>}
							/>
							<Route
								path={"/*"}
								element={<NotFound/>}
							/>
						</Routes>
					</main>
				</div>
			</ThemeProvider>
		</CacheProvider>
	)
}

export default App;
