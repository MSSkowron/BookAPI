alter table books 
add created_by bigint not null;

alter table books 
add constraint bookcreatedbyuserfk foreign key (created_by) references users(id);
