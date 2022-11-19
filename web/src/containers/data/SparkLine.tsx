import React, {useMemo} from "react";
import {styled, SxProps} from "@mui/material";
import {FF_METRIC_BASE_ZERO} from "../../App";

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

export type SparklineColours = "primary" | "secondary" | "inherit" | "success" | "warning" | "error";

interface Props {
	width?: number;
	height?: number;
	data: number[];
	color?: SparklineColours;
	sx?: SxProps;
	selectedIndex?: number;
	onSetSelectedIndex?: (v?: number) => void;
}

/**
 * Largely borrowed from https://github.com/Janpot/mui-plus/blob/master/packages/mui-plus/src/Sparkline/Sparkline.tsx
 * @licence MIT
 */
const SparkLine: React.FC<Props> = ({width = 150, height = 25, data, color = "primary", sx, selectedIndex, onSetSelectedIndex}): JSX.Element => {
	const margin = 3;
	const dotRadius = 3;

	const scaleLinear = (rangeMin: number, rangeMax: number, domainMin: number, domainMax: number): (x: number) => number => {
		const range = rangeMax - rangeMin;
		const domain = domainMax - domainMin;
		return (x: number): number => ((x - rangeMin) / range) * domain + domainMin;
	}

	const [scaleX, scaleY] = useMemo(() => {
		const min = FF_METRIC_BASE_ZERO ? Math.min(0, ...data) : Math.min(...data);
		const max = Math.max(...data);
		return [
			scaleLinear(0, data.length, margin, width - margin),
			scaleLinear(min, max, margin, height - margin),
		];
	}, [data, margin, width, height]);

	const selectedX = useMemo(() => {
		return scaleX(selectedIndex || 0);
	}, [selectedIndex]);

	const path = [];
	for (let i = 0; i < data.length; i++) {
		path.push(
			`${path.length <= 0 ? "M" : "L"} ${scaleX(i)} ${scaleY(data[i])}`
		);
	}

	return (
		<Root
			onMouseLeave={() => onSetSelectedIndex?.()}
			className={classes[`${color}Color` as keyof typeof classes]}
			transform="scale(1,-1)"
			sx={sx}
			width="100%"
			height="100%"
			viewBox={`0 0 ${width} ${height}`}>
			<path
				d={path.join(" ")}
				strokeWidth={1}
				stroke="currentcolor"
				fill="none"
			/>
			{selectedIndex != null && <path
				d={`M${selectedX} 0 L${selectedX} ${height}`}
				strokeWidth={2}
				stroke="currentcolor"
			/>}
			{data.map((d, i) => {
				const cx = scaleX(i) || 0;
				const cy = scaleY(d) || 0;
				const selected = selectedIndex != null && selectedX === cx;
				return <React.Fragment
					key={`${d}${i}`}>
					<path
						onMouseEnter={() => onSetSelectedIndex?.(i)}
						d={`M${cx} 0 L${cx} ${height}`}
						strokeWidth={10}
						stroke="transparent"
					/>
					{(selected || i === data.length - 1) && <circle
						onMouseEnter={() => onSetSelectedIndex?.(i)}
						cx={cx}
						cy={cy}
						r={selected ? dotRadius * 2 : dotRadius}
						fill="currentcolor"
					/>}
				</React.Fragment>
			})}
		</Root>
	);
}
export default SparkLine;
