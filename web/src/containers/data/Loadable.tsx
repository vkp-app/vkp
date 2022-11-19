import React, {useEffect, useState} from "react";
import {DEFAULT_LOAD_DELAY_MS} from "../../App";

interface Props {
	data: any;
	skeleton: any;
	loading: boolean;
	timeoutMs?: number;
}

const Loadable: React.FC<Props> = ({data, skeleton, loading, timeoutMs = DEFAULT_LOAD_DELAY_MS}): JSX.Element => {
	const [showLoading, setShowLoading] = useState<boolean>(false);
	const [duration, setDuration] = useState<number>(Date.now());

	useEffect(() => {
		if (loading) {
			setShowLoading(() => true);
			setDuration(() => Date.now());
			return;
		}
		// if the duration since we started loading
		// is bigger than the timeout, then don't
		// add any extra delay
		if ((Date.now() - duration) > timeoutMs) {
			setShowLoading(() => false);
			return;
		}
		setTimeout(() => {
			setShowLoading(() => false);
		}, timeoutMs);
	}, [loading]);

	return showLoading ? skeleton : data;
}
export default Loadable;
