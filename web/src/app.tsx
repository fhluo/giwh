import './app.css'
import {useEffect, useState} from 'react'
import {GetAssets, GetDefaultSharedWish, GetItems, GetLanguage, GetLocale} from '../wailsjs/go/main/App'
import {Tabs, TabsContent, TabsList, TabsTrigger} from './components/ui/tabs'
import {api} from '@/lib/models'
import {ScrollArea} from './components/ui/scroll-area'
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue} from './components/ui/select'
import {Avatar, AvatarFallback, AvatarImage} from './components/ui/avatar'
import {Label} from './components/ui/label'
import {Progress} from './components/ui/progress'
import {cn} from '@/lib/utils'
import {HoverCard, HoverCardContent, HoverCardTrigger} from './components/ui/hover-card'
import {Button} from './components/ui/button'
import {Loader2, RefreshCcw} from 'lucide-react'

function App() {
    let [uidList, setUIDList] = useState([] as string[])
    let [uid, setUID] = useState('')
    let [items, setItems] = useState([] as api.Item[])
    let [currentItems, setCurrentItems] = useState([] as api.Item[])

    let [wishes, setWishes] = useState([] as api.Wish[])
    let [sharedWishes, setSharedWishes] = useState([] as api.Wish[])
    let [currentSharedWishes, setCurrentSharedWishes] = useState([] as api.Wish[])
    let [currentWish, setCurrentWish] = useState('')

    let [locale, setLocale] = useState({} as api.Locale)
    let [locales, setLocales] = useState(new Map<string, api.Locale>())
    let [assets, setAssets] = useState({} as api.Assets)
    let [lang, setLang] = useState('')

    useEffect(() => {
        GetDefaultSharedWish().then(result => {
            setCurrentWish(result)
        })
        GetItems().then(result => {
            setItems(JSON.parse(result))
        })
        GetAssets().then(result => {
            setAssets(JSON.parse(result))
        })
        GetLanguage().then(result => {
            setLang(result)
        })
    }, [])


    const getLocale = async (lang: string): Promise<api.Locale> => {
        if (locales.has(lang)) {
            return locales.get(lang)!
        }
        let result: api.Locale = JSON.parse(await GetLocale(lang))
        setLocales(locales.set(lang, result))
        return result
    }

    useEffect(() => {
        getLocale(lang).then(result => {
            setLocale(result)
        })
    }, [lang])

    useEffect(() => {
        setUIDList([...new Set(items.map(item => item.uid))])
        setUID(uidList[0])

        let languages = new Set(items.map(item => item.lang))
        for (const lang of languages) {
            getLocale(lang).then()
        }

    }, [items])

    useEffect(() => {
        setSharedWishes(locale.sharedWishes)
        setWishes(locale.wishes)
    }, [locale])

    useEffect(() => {
        setCurrentItems(items.filter(item => {
            return item.uid == uid
        }))
    }, [items, uid])

    useEffect(() => {
        let types = new Set(currentItems.map(item => item.gacha_type))
        setCurrentSharedWishes(sharedWishes?.filter(wishType => types.has(wishType.key)))
    }, [currentItems])

    useEffect(() => {
        if (currentSharedWishes?.length > 0) {
            setCurrentWish(currentSharedWishes[0].key)
        }
    }, [currentSharedWishes])

    function getIcon(item: api.Item): string {
        const itemLocale = locales.get(item.lang)!
        if (item.name in itemLocale.charactersInverse) {
            return assets.characters[itemLocale.charactersInverse[item.name]]
        }
        if (item.name in itemLocale.weaponsInverse) {
            return assets.weapons[itemLocale.weaponsInverse[item.name]]
        }
        return ''
    }

    return (
        <div className="flex flex-col w-screen h-screen overflow-hidden">
            <div className="w-full px-4 pt-4 flex flex-row items-center justify-center">
                <div className="flex flex-col items-center mx-2">
                    <Avatar className="w-20 h-20 ring-offset-2 ring-2 cursor-pointer ring-[#ac7d3f] border-[#e9e5dd]">
                        <AvatarImage src={assets.characters ? assets.characters['Traveler'] : ''}
                                     className="bg-[#ac7d3f]"/>
                        <AvatarFallback>Traveler</AvatarFallback>
                    </Avatar>
                    <Select onValueChange={value => setUID(value)} defaultValue={uid}>
                        <SelectTrigger className="border-0 focus:ring-0 focus:ring-offset-0 items-end w-fit">
                            <SelectValue className="text-[#4d5562]">UID&nbsp;{uid}&nbsp;</SelectValue>
                        </SelectTrigger>
                        <SelectContent className="bg-[#e9e5dd] border-0">
                            {uidList.map(uid => <SelectItem value={uid} key={uid}
                                                            className="text-[#4d5562]">{uid}</SelectItem>)}
                        </SelectContent>
                    </Select>
                </div>
                <div className="flex flex-col space-y-3 mx-4">
                    {locale.sharedWishes?.filter(wishType => wishType.key != api.NoviceWishes).map(wishType => {
                        let pity = 90
                        let width = 'w-[270px]'
                        if (wishType.key == api.WeaponEventWish) {
                            pity = 80
                            width = 'w-[240px]'
                        }

                        let currentWishItems = currentItems.filter(item => {
                            return item.gacha_type == wishType.key || (
                                wishType.key == api.CharacterEventWish && item.gacha_type == api.CharacterEventWish2
                            )
                        }).sort((a, b) => {
                            if (a.id < b.id) {
                                return 1
                            } else if (a.id > b.id) {
                                return -1
                            } else {
                                return 0
                            }
                        })

                        let value = currentWishItems.findIndex(item => item.rank_type == api.FiveStar)

                        return (
                            <div className="grid grid-cols-3">
                                <Label>{wishes?.find(w => w.key == wishType.key)!.name}</Label>
                                <Progress value={100 * value / pity} max={pity}
                                          className={cn('bg-slate-200', width, 'col-span-2')}></Progress>
                            </div>
                        )
                    })}
                </div>
                <Button variant="subtle" className="bg-slate-200 rounded-full h-12 mx-4">
                    {/*"animate-spin"*/}
                    <RefreshCcw className={cn("mr-2 h-4 w-4")}/>
                    Update
                </Button>
            </div>
            <Tabs defaultValue={api.CharacterEventWishAndCharacterEventWish2} onValueChange={value => {
                let wishType = sharedWishes?.find(v => v.key == value)
                setCurrentWish(wishType!.key)
            }} className="flex flex-col w-full h-full">
                <TabsList className="bg-slate-200">
                    {
                        currentSharedWishes?.map(wishType => {
                                return (
                                    <TabsTrigger value={wishType.key} key={wishType.key}>
                                        {wishes.find(w => w.key == wishType.key)!.name}
                                    </TabsTrigger>
                                )
                            },
                        )
                    }
                </TabsList>
                {
                    currentSharedWishes?.map(wishType => {
                        let currentWishItems = currentItems?.filter(item => {
                            return item.gacha_type == wishType.key || (
                                wishType.key == api.CharacterEventWish && item.gacha_type == api.CharacterEventWish2
                            )
                        }).sort((a, b) => {
                            if (a.id < b.id) {
                                return -1
                            } else if (a.id > b.id) {
                                return 1
                            } else {
                                return 0
                            }
                        })

                        // console.log(items)
                        let items_ = currentWishItems.filter(item => {
                            return item.rank_type == api.FiveStar
                        })
                        return (
                            <TabsContent value={wishType.key.toString()}
                                         className="w-full h-full border-0 p-0"
                                         key={wishType.key}>
                                <div className="w-full h-full flex flex-col">
                                    <ScrollArea
                                        className="h-72 grow px-6 py-2">
                                        <div
                                            className="flex flex-row flex-wrap items-center justify-center w-full gap-x-5 gap-y-5 transition duration-200 p-2">
                                            {
                                                items_.map((item, index) => {
                                                        // <Item item={item} pulls={0} key={item.id} } wishes={wishes}/>
                                                        let icon = getIcon(item)

                                                        let pulls = 0
                                                        if (index == 0) {
                                                            pulls = currentWishItems.indexOf(item) + 1
                                                        } else {
                                                            pulls = currentWishItems.indexOf(item) - currentWishItems.indexOf(items_[index - 1])
                                                        }
                                                        return (
                                                            <HoverCard>
                                                                <HoverCardTrigger
                                                                    className=" flex bg-[#e9e5dd] flex-col w-fit items-center select-none transition duration-300 hover:ring-2 hover:ring-offset-2 shadow rounded-lg">
                                                                    <div>
                                                                        <div className="w-20">
                                                                            <img alt={item.name}
                                                                                 className={cn('pointer-events-none rounded-t-lg shadow-inner', 'bg-[#ac7d3f]')}
                                                                                 src={icon}/>
                                                                        </div>
                                                                    </div>
                                                                    <div
                                                                        className="w-20 text-sm text-[#4d5562] border-t text-center rounded-b-lg leading-normal tracking-wider">{pulls}</div>
                                                                </HoverCardTrigger>
                                                                <HoverCardContent
                                                                    className="flex flex-col space-y-3 items-center text-sm text-[#4d5562] bg-[#e9e5dd] border-0 w-fit">
                                                                    <Label className="font-semibold">{item.name}</Label>
                                                                    <Label>{item.item_type}</Label>
                                                                    <Label>{wishes?.find(wishType => wishType.key == item.gacha_type)!.name}</Label>
                                                                    <Label>{item.time}</Label>
                                                                </HoverCardContent>
                                                            </HoverCard>
                                                        )
                                                    },
                                                )
                                            }
                                        </div>
                                    </ScrollArea>
                                </div>
                            </TabsContent>
                        )
                    })
                }
            </Tabs>

        </div>
    )
}

export default App
