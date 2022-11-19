import React from "react";
import {Box} from "@mui/material";
import {Navigate} from "react-router-dom";
import StandardLayout from "../layout/StandardLayout";

const Home: React.FC = (): JSX.Element => {
	return <StandardLayout>
		<Box
			sx={{p: 2}}>
			<Navigate to="/tenants"/>
		</Box>
	</StandardLayout>
}
export default Home;
