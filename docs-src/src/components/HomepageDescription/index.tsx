import type { ReactNode } from "react";
import styles from "./styles.module.css";

export default function HomepageFeatures(): ReactNode {
    return (
        <section className={styles.features}>
            <div className="container">
                <h1>Welcome to PyramID-Engine</h1>
                <p>Unlock the power of dynamic agent access management and seamless integration with PyramID-Engine. Our solution is designed to enhance your organization's adoption of agents efficiently and securely, providing robust features that cater to your unique needs.</p>
            </div >
            <div className="container" style={{ paddingTop: "2rem" }}>
                <h3>How it Works</h3>
                <ol>
                    <li><b>Connect IDPs:</b> Integrate with identity providers to extract user group memberships.</li>
                    <li><b>Bind Sources to Groups:</b> Assign data sources or parts of sources to specific user groups for role-based access.</li>
                    <li><b>Dynamic Tool Access:</b> Use the PyramID-Engine SDK to dynamically load agent tools based on user permissions.</li>
                    <li><b>Test and Validate:</b> Utilize our playground with reference agents to ensure configurations work as expected.</li>
                </ol>
            </div >
            <div className="container" style={{ paddingTop: "2rem" }}>
                <h3>Use Cases</h3>
                <ul>
                    <li>Enhanced Security: Prevent unauthorized access by dynamically adjusting agent tool availability.</li>
                    <li>Operational Efficiency: Reduce unnecessary operations and save resources by aligning agent tool access with user permissions.</li>
                    <li>Tailored Solutions: Customize agent tool functionalities to fit your organization's unique needs.</li>
                </ul>
            </div >
            <div className="container" style={{ paddingTop: "2rem" }}>
                <h3>Under the Hood</h3>
                <p>Explore the workings of PyramID-Engine with our simplified architecture diagram. This visual representation showcases:</p>
                <ul>
                    <li><b>IDP Integration:</b> How PyramID-Engine connects with various identity providers to extract and manage user group memberships.</li>
                    <li><b>Source Binding:</b> The process of linking data sources to specific user groups, ensuring role-based access control.</li>
                    <li><b>Dynamic Tool Loading:</b> The mechanism that enables agents to load tools dynamically based on user permissions, ensuring they have access to only authorized functionalities.</li>
                    <li><b>Data Flow:</b> The seamless flow of data between users, agents, and connected sources, facilitated by the PyramID-Engine backend.</li>
                </ul>
                <img style={{ width: "80%", display: "block", margin: "0 auto" }} src={require("@site/static/img/homepage/PyramID-Engine-Overview.drawio.png").default} alt="Architecture Diagram" />
            </div>
            <div className="container" style={{ paddingTop: "2rem", paddingBottom: "2rem" }}>
                <h3>Get Started</h3>
                <p>Ready to experience the benefits of PyramID-Engine? Check out our <a href="/docs/intro">quick tutorial</a> to get started in just 5 minutes!</p>
            </div >
        </section >
    );
}