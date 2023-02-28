import ICTSCNavBar from "../components/Navbar";
import MarkdownPreview from "../components/MarkdownPreview";
import ICTSCCard from "../components/Card";
import Head from "next/head";
import { rule, site } from "../components/_const";

const Home = () => {
  return (
    <>
      <Head>
        <title>ルール - {site}</title>
      </Head>
      <ICTSCNavBar />
      <h1 className={"title-ictsc text-center py-12"}>ルール</h1>
      <div className={"container-ictsc"}>
        <ICTSCCard>
          <MarkdownPreview content={rule?.replace(/\\n/g, "\n") ?? ""} />
        </ICTSCCard>
      </div>
    </>
  );
};

export default Home;
