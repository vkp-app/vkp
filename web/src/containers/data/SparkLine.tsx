import React, {useMemo} from "react";
import {styled, SxProps} from "@mui/material";

const classes = {
	primaryColor: "primaryColor",
	secondaryColor: "secondaryColor",
	successColor: "successColor",
	warningColor: "warningColor",
	errorColor: "errorColor",
	inheritColor: "inheritColor",
} as const;

const Root = styled("svg")(({theme}) => ({
	color: theme.palette.primary.main,
	[`&.${classes.secondaryColor}`]: {
		color: theme.palette.secondary.main,
	},
	[`&.${classes.successColor}`]: {
		color: theme.palette.success.main,
	},
	[`&.${classes.warningColor}`]: {
		color: theme.palette.warning.main,
	},
	[`&.${classes.errorColor}`]: {
		color: theme.palette.error.main,
	},
	[`&.${classes.inheritColor}`]: {
		color: "inherit",
	},
}));

interface Props {
	width?: number;
	height?: number;
	data: number[];
	color?: "primary" | "secondary" | "inherit" | "success" | "warning" | "error";
	sx?: SxProps;
}

const SparkLine: React.FC<Props> = ({width = 150, height = 25, data, color = "primary", sx}): JSX.Element => {
	const margin = 3;
	const dotRadius = 3;

	const scaleLinear = (rangeMin: number, rangeMax: number, domainMin: number, domainMax: number): (x: number) => number => {
		const range = rangeMax - rangeMin;
		const domain = domainMax - domainMin;
		return (x: number): number => ((x - rangeMin) / range) * domain + domainMin;
	}

	const [scaleX, scaleY] = useMemo(() => {
		const min = Math.min(...data);
		const max = Math.max(...data);
		return [
			scaleLinear(0, data.length, margin, width - margin),
			scaleLinear(min, max, margin, height - margin),
		];
	}, [data, margin, width, height]);

	const path = [];
	for (let i = 0; i < data.length; i++) {
		path.push(
			`${path.length <= 0 ? "M" : "L"} ${scaleX(i)} ${scaleY(data[i])}`
		);
	}

	const cx = scaleX(data.length - 1) || 0;
	const cy = scaleY(data[data.length - 1]) || 0;

	return (
		<Root
			className={classes[`${color}Color` as keyof typeof classes]}
			transform="scale(1,-1)"
			sx={sx}
			width={width}
			height={height}
			viewBox={`0 0 ${width} ${height}`}>
			<path
				d={path.join(" ")}
				strokeWidth={1}
				stroke="currentcolor"
				fill="none"
			/>
			<circle cx={cx} cy={cy} r={dotRadius} fill="currentcolor"></circle>
		</Root>
	);
}
export default SparkLine;
