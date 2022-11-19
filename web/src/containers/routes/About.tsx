import {Avatar, Card, CardHeader} from "@mui/material";
import React from "react";
import StandardLayout from "../layout/StandardLayout";

const About: React.FC = (): JSX.Element => {
	return <StandardLayout>
		<Card
			sx={{p: 2}}>
			<CardHeader
				sx={{p: 0, pl: 1}}
				title="About"
				titleTypographyProps={{fontFamily: "Figtree"}}
			/>
			<CardHeader
				title="Virtual Kubernetes Platform v0.0.0"
				subheader="Management platform for Virtual Kubernetes Clusters."
				avatar={<Avatar
					sx={{width: 64, height: 64}}
					src="/src/img/kubernetes-icon-color.svg"
					alt="VKP icon"
				/>}
			/>
		</Card>
	</StandardLayout>
}
export default About;
