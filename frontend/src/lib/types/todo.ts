export type Priority = 'LOW' | 'MEDIUM' | 'HIGH' | 'URGENT';

export type Task = {
	Id: string;
	Title: string;
	Description: string;
	DueDate: string;
	IsDone: boolean;
	Priority: Priority;
	Category: string;
	DueWarning: string;
	ChildTasks: Task[];
	ParentId?: string;
};

export type TaskResponse = {
	id: string;
	title: string;
	description: string;
	dueDate: string;
	status: string;
	priority: string;
	category: string;
	dueWarning: string;
	childTasks: TaskResponse[];
	parentId?: string;
	createdAt: string;
	updatedAt: string;
};

export type CreateTaskRequest = {
	title: string;
	description: string;
	dueDate: string;
	priority: string;
	category: string;
	parentId?: string;
};
