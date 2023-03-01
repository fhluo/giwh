import {GetUIDs} from '../../wailsjs/go/main/App'
import {useEffect, useState} from 'react'
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue} from './ui/select'
import {Avatar} from './ui/avatar'

interface SidebarProps {
    currentUID: number;
}

function Sidebar({currentUID}: SidebarProps) {
    let [uidList, setUIDList] = useState([] as number[])

    useEffect(() => {
        GetUIDs().then(result => {
            setUIDList(result)
        })
    })

    return (
        <div>
            <div>
                <Avatar>

                </Avatar>
            </div>
            <div className="">
                <Select>
                    <SelectTrigger>
                        <SelectValue placeholder="UID">{currentUID}</SelectValue>
                    </SelectTrigger>
                    <SelectContent>
                        {uidList.map(uid => <SelectItem value={uid.toString()}>{uid}</SelectItem>)}
                    </SelectContent>
                </Select>
            </div>
        </div>
    )
}

export default Sidebar