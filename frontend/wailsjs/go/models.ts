export namespace models {
	
	export class Education {
	    id: number;
	    profile_id: number;
	    institution: string;
	    course: string;
	    start_date: string;
	    end_date: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new Education(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.profile_id = source["profile_id"];
	        this.institution = source["institution"];
	        this.course = source["course"];
	        this.start_date = source["start_date"];
	        this.end_date = source["end_date"];
	        this.description = source["description"];
	    }
	}
	export class Experience {
	    id: number;
	    profile_id: number;
	    company: string;
	    position: string;
	    start_date: string;
	    end_date: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new Experience(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.profile_id = source["profile_id"];
	        this.company = source["company"];
	        this.position = source["position"];
	        this.start_date = source["start_date"];
	        this.end_date = source["end_date"];
	        this.description = source["description"];
	    }
	}
	export class Profile {
	    id: number;
	    first_name: string;
	    last_name: string;
	    email: string;
	    phone: string;
	    address: string;
	    age: number;
	    photo: string;
	    objective: string;
	    linkedin: string;
	    github: string;
	    website: string;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Profile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.first_name = source["first_name"];
	        this.last_name = source["last_name"];
	        this.email = source["email"];
	        this.phone = source["phone"];
	        this.address = source["address"];
	        this.age = source["age"];
	        this.photo = source["photo"];
	        this.objective = source["objective"];
	        this.linkedin = source["linkedin"];
	        this.github = source["github"];
	        this.website = source["website"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	export class Project {
	    id: number;
	    profile_id: number;
	    name: string;
	    description: string;
	    url: string;
	
	    static createFrom(source: any = {}) {
	        return new Project(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.profile_id = source["profile_id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.url = source["url"];
	    }
	}
	export class Settings {
	    id: number;
	    language: string;
	    template: string;
	    labels: Record<string, string>;
	    show_photo: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.language = source["language"];
	        this.template = source["template"];
	        this.labels = source["labels"];
	        this.show_photo = source["show_photo"];
	    }
	}

}

export namespace parser {
	
	export class ExtractedData {
	    Profile: models.Profile;
	    Education: models.Education[];
	    Experience: models.Experience[];
	    Projects: models.Project[];
	
	    static createFrom(source: any = {}) {
	        return new ExtractedData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Profile = this.convertValues(source["Profile"], models.Profile);
	        this.Education = this.convertValues(source["Education"], models.Education);
	        this.Experience = this.convertValues(source["Experience"], models.Experience);
	        this.Projects = this.convertValues(source["Projects"], models.Project);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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

