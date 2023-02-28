import {GetUIDs} from '../../wailsjs/go/main/App.js'
import {useEffect, useState} from 'react'

interface SideBarProps {
    currentUID: number;
}

function SideBar({currentUID}: SideBarProps) {
    let [uidList, setUIDList] = useState([] as number[])

    useEffect(() => {
        GetUIDs().then(result => {
            setUIDList(result)
        })
    })

    return (
        <div>
            <div className="flex flex-row shadow ring-2 ring-cyan-500 rounded items-center bg-white/25 text-sm">
                <label
                    className="inline-block font-bold px-4 py-1 border-r bg-white/75 rounded-l-lg select-none">UID</label>
                <button
                    className="px-4 py-1 cursor-pointer hover:bg-gray-300/50 hover:shadow hover:-top-1 select-none">{currentUID}</button>
                {
                    uidList.map(uid =>
                        <div
                            className="px-5 py-2 cursor-pointer hover:bg-gray-300/50 hover:shadow hover:-top-1 select-none"
                            onClick={() => currentUID = uid}>{uid}
                        </div>
                    )
                }
            </div>
        </div>
    )
}

export default SideBar