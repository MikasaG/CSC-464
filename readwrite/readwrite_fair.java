package readwrite;

import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.Semaphore;


public class readwrite_fair {
	static Semaphore mutex = new Semaphore(1);
	static Semaphore roomEmpty = new Semaphore(1);
	static int readers = 0;
	static final int reader_num = 1000;
	static final int writer_num = 1000;
	
	private static void read(String s) {
		System.out.println(s+" is reading");
	}
	
	private static void write(String s) {
		System.out.println(s+" is writing");
	}
	
	static class reader implements Runnable{
		int id;

		public void setId(int id) {
			this.id = id;
		}
		
		@Override
		public void run()  {
			try {
				mutex.acquire();
				readers++;
				if (readers == 1) {
					roomEmpty.acquire();
				}
				mutex.release();
				read("Reader " + id);
				
				mutex.acquire();
				readers--;
				if (readers == 0) {
					roomEmpty.release();
				}
				mutex.release();
			} catch (InterruptedException e){
				e.printStackTrace();
			}
	    }
	}
	
	static class writer implements Runnable{
		int id;

		public void setId(int id) {
			this.id = id;
		}
		
		@Override
		public void run()  {
			try {
				roomEmpty.acquire();
				
				write("Writer "+id);
				
				roomEmpty.release();
			} catch (InterruptedException e){
				e.printStackTrace();
			}
	    }
	}

	
	public static void main (String[] args) {
		Thread[] threads = new Thread[reader_num + writer_num];
		for(int i = 0; i < reader_num; i++){
			reader r = new reader();
			r.setId(i);
			threads[i] = new Thread(r);
		}
		for(int i = reader_num; i < threads.length; i++){
			writer w = new writer();
			w.setId(i);
			threads[i] = new Thread(w);
		}
		List<Thread> l = Arrays.asList(threads);
		Collections.shuffle(l);
		long startTime = System.currentTimeMillis(); 
		for(int i = 0; i < l.size(); i++){
				  l.get(i).start();
		}
		for(int j = 0; j < l.size(); j++){
			  try {
				  l.get(j).join();
			  } catch (Exception e) {
				  e.printStackTrace();
			  } 
		}
		long endTime=System.currentTimeMillis(); //get end time
		System.out.println("Time Cost "+(endTime-startTime)+"ms");	
	}
}
