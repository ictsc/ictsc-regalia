import React from "react";

const ICTSCNavBar: React.FC = () => {
  return (
      <div className={"navbar bg-primary text-primary-content"}>
        <div className={'flex-1'}>
          <a className="btn btn-ghost normal-case text-xl">ICTSC</a>
        </div>
        <div className={'flex-none'}>
          <ul className="menu menu-horizontal p-0">
            <li><a>問題</a></li>
          </ul>
        </div>
      </div>
  )
}

export default ICTSCNavBar;