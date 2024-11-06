from django.core.management import BaseCommand
from django.conf import settings
from kombu import Exchange, Queue

from core.rabbitmq import create_rabbitmq_connection
from core.service import create_video_service_factory

class Command(BaseCommand):
    help = 'Register processed video path'
    
    def handle(self, *args, **options):
        confirmation_exchange = settings.CONFIRMATION_EXCHANGE
        confirmation_queue = settings.CONFIRMATION_QUEUE
        confirmation_key = settings.CONFIRMATION_KEY
       
        self.stdout.write(self.style.SUCCESS('Starting consumer...'))
       
        exchange = Exchange(confirmation_exchange, type='direct', auto_delete= True)
        queue = Queue(confirmation_queue, exchange, confirmation_key)
       
        with create_rabbitmq_connection() as conn:
            with conn.Consumer(queue, callbacks=[self.process_message]):
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
        