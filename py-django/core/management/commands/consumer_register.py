from django.core.management import BaseCommand
from decouple import config
from kombu import Exchange, Queue

from core.rabbitmq import create_rabbitmq_connection
from core.service import create_video_service_factory

class Command(BaseCommand):
    help = 'Register processed video path'
    
    def handle(self, *args, **options):
       confirmationExchange = config('CONFIRMATION_EXCHANGE')
       confirmationQueue = config('CONFIRMATION_QUEUE')
       confirmationKey = config('CONFIRMATION_KEY') 
       
       self.stdout.write(self.style.SUCCESS('Starting consumer...'))
       
       exchange = Exchange(confirmationExchange, type='direct', auto_delete= True)
       queue = Queue(confirmationQueue, exchange, confirmationKey)
       
       with create_rabbitmq_connection() as conn:
           with conn.Consumer(queue, callback=[self.process_message]):
               while True:
                   self.stdout.write(self.style.SUCCESS('Waiting for messages...'))
                   conn.drain_events()
           
    def process_message(self, body, message):
        self.stdout.write(self.style.SUCCESS(f'Processing message: {body}'))
        create_video_service_factory().register_processed_video_path(
            video_id=body['video_id'], 
            video_path=body['path']
        )
        message.ack()
        