const ProblemCard = () => {
  return (
      <a className={'border p-4 hover:bg-base-200 hover:cursor-pointer rounded-md shadow-sm min-h-[212px] justify-between flex flex-col'}>
        <div>
          <span className={'font-bold text-2xl text-primary pr-2'}>1</span>
          <span className={'text-xl font-bold'}>問題タイトル</span>
        </div>
        <div>
          <div className={'text-right'}>100/100pt</div>
          <div className={'font-bold text-primary'}>問題文へ→</div>
        </div>
      </a>
  )
}

export default ProblemCard