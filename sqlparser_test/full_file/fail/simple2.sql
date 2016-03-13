/*
 * This is a SQL file that will have a syntax error. Multiple infact.
 * However I only expect the parser to find the first one. And stop.
 */
 
 select 1,b form b where a.b == 1;
 
 -- The next one I would not expect to be found.
 
 select a,b from c where a.id = 1;