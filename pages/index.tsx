import CommonLayout from "@/layouts/CommonLayout";
import { rule } from "@/components/_const";
import MarkdownPreview from "@/components/MarkdownPreview";
import ICTSCCard from "@/components/Card";

const Home = () => (
  <CommonLayout title={"ルール"}>
    <div className={"container-ictsc"}>
      <ICTSCCard>
        <MarkdownPreview content={rule?.replace(/\\n/g, "\n") ?? ""} />
      </ICTSCCard>
    </div>
  </CommonLayout>
);

export default Home;
