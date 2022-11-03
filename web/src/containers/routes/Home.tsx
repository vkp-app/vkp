import React from "react";
import StandardLayout from "../layout/StandardLayout";
import {Box, Link as MuiLink, Typography} from "@mui/material";
import {Link} from "react-router-dom";

const Home: React.FC = (): JSX.Element => {
	return <StandardLayout>
		<Box
			sx={{p: 2}}>
			<Typography>
				This page is a dummy that will be overwritten by a branding page.
				<br/>
				You&apos;re probably looking for&nbsp;
				<MuiLink
					component={Link}
					to={"/clusters/tenant-sample"}>
					/clusters
				</MuiLink>
			</Typography>
		</Box>
	</StandardLayout>
}
export default Home;
