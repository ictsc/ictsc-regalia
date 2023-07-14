import LoadingSpinner from "./LoadingSpinner";

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
