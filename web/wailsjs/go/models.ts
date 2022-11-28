export namespace api {
	
	export class Item {
	    count: number;
	    gacha_type: number;
	    id: number;
	    item_id: string;
	    item_type: string;
	    lang: string;
	    name: string;
	    rank_type: number;
	    // Go type: Time
	    time: any;
	    uid: number;
	
	    static createFrom(source: any = {}) {
	        return new Item(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.count = source["count"];
	        this.gacha_type = source["gacha_type"];
	        this.id = source["id"];
	        this.item_id = source["item_id"];
	        this.item_type = source["item_type"];
	        this.lang = source["lang"];
	        this.name = source["name"];
	        this.rank_type = source["rank_type"];
	        this.time = this.convertValues(source["time"], null);
	        this.uid = source["uid"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace repository {
	
	export class Item {
	    count: number;
	    gacha_type: number;
	    id: number;
	    item_id: string;
	    item_type: string;
	    lang: string;
	    name: string;
	    rank_type: number;
	    // Go type: api.Time
	    time: any;
	    uid: number;
	    pulls: number;
	    icon: string;
	
	    static createFrom(source: any = {}) {
	        return new Item(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.count = source["count"];
	        this.gacha_type = source["gacha_type"];
	        this.id = source["id"];
	        this.item_id = source["item_id"];
	        this.item_type = source["item_type"];
	        this.lang = source["lang"];
	        this.name = source["name"];
	        this.rank_type = source["rank_type"];
	        this.time = this.convertValues(source["time"], null);
	        this.uid = source["uid"];
	        this.pulls = source["pulls"];
	        this.icon = source["icon"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

