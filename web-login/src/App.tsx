import React, {PropsWithChildren, useEffect, useMemo} from "react";
import createCache from "@emotion/cache";
import {createTheme, CssBaseline, useMediaQuery} from "@mui/material";
import {CacheProvider} from "@emotion/react";
import {ThemeProvider} from "@mui/styles";
import {makeStyles} from "tss-react/mui";
import "typeface-roboto";

const useStyles = makeStyles()(() => ({
	root: {
		display: "flex",
	},
	content: {
		flexGrow: 1
	},
}));


const App: React.FC<PropsWithChildren> = ({children}): JSX.Element => {
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
						variant: "elevation"
					}
				},
				MuiCard: {
					defaultProps: {
						variant: "elevation"
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
							variant: "elevation",
							elevation: 4,
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
					<main className={classes.content}>
						{children}
					</main>
				</div>
			</ThemeProvider>
		</CacheProvider>
	);
}

export default App;
