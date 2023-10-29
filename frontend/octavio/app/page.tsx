import { rule } from "@/components/_const";
import ICTSCCard from "@/components/card";
import MarkdownPreview from "@/components/markdown-preview";
import ICTSCTitle from "@/components/title";

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
