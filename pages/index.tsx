import BaseLayout from "@/layouts/BaseLayout";
import { rule } from "@/components/_const";
import MarkdownPreview from "@/components/MarkdownPreview";
import ICTSCCard from "@/components/Card";

const Home = () => (
  <BaseLayout title={"ルール"}>
    <h1 className={"title-ictsc text-center py-12"}>ルール</h1>
    <div className={"container-ictsc"}>
      <ICTSCCard>
        <MarkdownPreview content={rule?.replace(/\\n/g, "\n") ?? ""} />
      </ICTSCCard>
    </div>
  </BaseLayout>
);

export default Home;
