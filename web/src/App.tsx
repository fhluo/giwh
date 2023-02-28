import './App.css'
import {api} from './models'
import {useEffect, useState} from 'react'
import {Get5Stars, GetUIDs, GetWishName} from '../wailsjs/go/main/App'
import SideBar from './lib/SideBar'
import {repository} from '../wailsjs/go/models'
import Avatar from './lib/Avatar'

function App() {
    let [uidList, setUIDList] = useState([] as number[])
    let [currentUID, setCurrentUID] = useState(0)
    let [items, setItems] = useState([] as repository.Item[])

    let currentWishType = api.SCharacterEventWish


    useEffect(() => {
        GetUIDs().then(result => {
            setUIDList(result)
            if (uidList.length != 0) {
                setCurrentUID(uidList[1])
            }
        })
    })

    useEffect(() => {
        Get5Stars(currentUID, currentWishType).then(result => {
            setItems(result)
        })
    }, [currentUID, currentWishType])

    return (
        <div className="flex flex-row my-3 mx-3">
            <SideBar currentUID={currentUID}/>
            <div className="flex flex-col px-5 w-full h-screen">
                <div className="flex flex-row flex-wrap  items-center w-fit  mx-12">
                    {
                        api.SharedWishTypes.map(wishType => {
                                let [wishName, setWishName] = useState('')
                                useEffect(() => {
                                    GetWishName(wishType).then(result => {
                                        setWishName(result)
                                    })
                                })
                                return (
                                    <div
                                        className="option transition duration-200 px-8 py-2 cursor-pointer hover:bg-gray-300/25 select-none leading-relaxed tracking-wider {currentWishType===wishType?'bg-gray-300/25 border-b-2 border-b-blue-500':''}"
                                        onClick={() => currentWishType = wishType}>
                                        {wishName}
                                    </div>
                                )
                            }
                        )
                    }
                </div>

                <div className=" py-12 px-12 tracking-wider space-y-3 w-full">
                    <div className="flex flex-row flex-wrap w-full gap-x-6 gap-y-6 transition duration-200">
                        {
                            items.map(item => <Avatar icon={item.icon} name={item.name} pulls={item.pulls}/>)
                        }
                    </div>
                </div>
            </div>
        </div>
    )
}

export default App
