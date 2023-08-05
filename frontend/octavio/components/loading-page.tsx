import LoadingSpinner from "./loading-spinner";

function LoadingPage() {
  return (
    <div
      className="flex justify-center items-center h-screen"
      data-testid="loading"
    >
      <LoadingSpinner />
    </div>
  );
}

export default LoadingPage;
