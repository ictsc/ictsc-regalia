import ICTSCCard from "@/components/Card";
import MarkdownPreview from "@/components/MarkdownPreview";
import { rule } from "@/components/_const";
import CommonLayout from "@/layouts/CommonLayout";

function Home() {
  return (
    <CommonLayout title="ルール">
      <div className="container-ictsc">
        <ICTSCCard>
          <MarkdownPreview content={rule} />
        </ICTSCCard>
      </div>
    </CommonLayout>
  );
}

export default Home;
