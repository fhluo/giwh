import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue} from './ui/select'
import {Avatar} from './ui/avatar'
import React from 'react'

interface SidebarProps {
    uid: string;
    setUID: React.Dispatch<React.SetStateAction<string>>;
    uidList: string[];
}

function Sidebar({uid, setUID, uidList}: SidebarProps) {
    return (
        <div className="w-1/5">
            <div>
                <Avatar>

                </Avatar>
            </div>
            <div className="">
                <Select onValueChange={value => setUID(value)}>
                    <SelectTrigger className="bg-white">
                        <SelectValue placeholder="UID">{uid}</SelectValue>
                    </SelectTrigger>
                    <SelectContent>
                        {uidList.map(uid => <SelectItem value={uid} key={uid}>{uid}</SelectItem>)}
                    </SelectContent>
                </Select>
            </div>
        </div>
    )
}

export default Sidebar