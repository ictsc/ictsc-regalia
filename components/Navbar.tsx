import Link from "next/link";

const ICTSCNavBar = () => {
  return (
      <div className={"navbar bg-primary text-primary-content"}>
        <div className={'flex-1'}>
          <Link href={'/'} className="btn btn-ghost normal-case text-xl">ICTSC</Link>
        </div>
        <div className={'flex-none'}>
          <ul className="menu menu-horizontal p-0">
            <li>
              <Link href={'/'}>ルール</Link>
            </li>
            <li>
              <Link href={'/problems'}>問題</Link>
            </li>
          </ul>
        </div>
      </div>
  )
}

export default ICTSCNavBar;