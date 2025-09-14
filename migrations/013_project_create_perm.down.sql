-- Remove global permission key for creating projects
delete from permissions where key = 'project.create';
