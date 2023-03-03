import {cn} from '@/lib/utils'
import {api} from '@/lib/models'

interface AvatarProps {
    item: api.Item;
    icon: string;
    pulls: number;
}

function Item({item, icon, pulls}: AvatarProps) {
    return (
        <div
            className="cursor-pointer flex flex-col w-fit items-center select-none shadow rounded-lg hover:bg-white/25 transition duration-200">
            <div>
                <div className="w-20">
                    <img alt={item.name} className={cn('pointer-events-none rounded-t-lg shadow-inner')}
                         src={icon} style={{backgroundColor: '#ac7d3f'}}/>
                </div>
            </div>
            {/*<div className="font-semibold text-gray-900 tracking-wide leading-relaxed">{name}</div>*/}
            <div style={{backgroundColor: '#e9e5dd', color: '#4d5562'}}
                 className="text-sm border-t w-full text-center rounded-b-lg leading-relaxed tracking-wider">{pulls}</div>
        </div>
    )
}

export default Item