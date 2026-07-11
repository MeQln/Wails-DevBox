export namespace main {
	
	export class DialogFilter {
	    name: string;
	    extensions: string[];
	
	    static createFrom(source: any = {}) {
	        return new DialogFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.extensions = source["extensions"];
	    }
	}
	export class HashResult {
	    size: number;
	    md5: string;
	    sha1: string;
	    sha256: string;
	    sha384: string;
	    sha512: string;
	
	    static createFrom(source: any = {}) {
	        return new HashResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.size = source["size"];
	        this.md5 = source["md5"];
	        this.sha1 = source["sha1"];
	        this.sha256 = source["sha256"];
	        this.sha384 = source["sha384"];
	        this.sha512 = source["sha512"];
	    }
	}
	export class ImageInfo {
	    width: number;
	    height: number;
	    format: string;
	    size_bytes: number;
	    data_base64: string;
	
	    static createFrom(source: any = {}) {
	        return new ImageInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.width = source["width"];
	        this.height = source["height"];
	        this.format = source["format"];
	        this.size_bytes = source["size_bytes"];
	        this.data_base64 = source["data_base64"];
	    }
	}
	export class PasswordOptions {
	    length: number;
	    upper: boolean;
	    lower: boolean;
	    digit: boolean;
	    symbol: boolean;
	    excludeAmbiguous: boolean;
	
	    static createFrom(source: any = {}) {
	        return new PasswordOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.length = source["length"];
	        this.upper = source["upper"];
	        this.lower = source["lower"];
	        this.digit = source["digit"];
	        this.symbol = source["symbol"];
	        this.excludeAmbiguous = source["excludeAmbiguous"];
	    }
	}
	export class PortCheckResult {
	    host: string;
	    port: number;
	    open: boolean;
	    latency_ms: number;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new PortCheckResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.host = source["host"];
	        this.port = source["port"];
	        this.open = source["open"];
	        this.latency_ms = source["latency_ms"];
	        this.message = source["message"];
	    }
	}
	export class PortEntry {
	    port: number;
	    pid: number;
	    process_name: string;
	    address: string;
	
	    static createFrom(source: any = {}) {
	        return new PortEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.port = source["port"];
	        this.pid = source["pid"];
	        this.process_name = source["process_name"];
	        this.address = source["address"];
	    }
	}

}

