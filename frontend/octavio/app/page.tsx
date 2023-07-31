import ICTSCCard from "@/components/Card";
import MarkdownPreview from "@/components/MarkdownPreview";
import ICTSCTitle from "@/components/Title";
import { rule } from "@/components/_const";

function Home() {
  return (
    <>
      <ICTSCTitle title="ルール" />
      <main className="container-ictsc">
        <ICTSCCard>
          <MarkdownPreview content={rule} />
        </ICTSCCard>
      </main>
    </>
  );
}

export default Home;
