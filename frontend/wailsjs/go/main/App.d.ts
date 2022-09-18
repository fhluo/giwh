// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {api} from '../models';
import {repository} from '../models';

export function Get5Stars(arg1:number,arg2:api.SharedWishType):Promise<Array<repository.Item>>;

export function GetPity(arg1:api.Rarity,arg2:api.SharedWishType):Promise<number>;

export function GetProgress(arg1:number,arg2:api.SharedWishType):Promise<number>;

export function GetSharedWishName(arg1:api.SharedWishType):Promise<string>;

export function GetSharedWishTypes():Promise<Array<api.SharedWishType>>;

export function GetUIDs():Promise<Array<number>>;
